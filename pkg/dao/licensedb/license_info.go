package licensedb

import (
	"fmt"
	"lisence/pkg/dao/dbpool"
	"lisence/pkg/models"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// LicenseInfoModule license_infos表的数据对象LicenseInfo操作类
type LicenseInfoModule struct{}

// ApplyStatus 用户查询license是否申请成功返回的信息(管理员是否审批，license，license过期时间)
type ApplyStatus struct {
	Status      bool      `json:"status"`
	License     string    `json:"license"`
	ExpiredTime time.Time `json:"expired_time"`
}

// UserApply 用户申请License时需要提交的信息
type UserApply struct {
	Quota       uint16 `json:"quota"`       // 配额
	Period      uint16 `json:"period"`      // 申请的服务期限
	Uuid        string `json:"uuid"`        // 该用户的唯一标识符
	Corporation string `json:"corporation"` // 公司名称
	Version     string `json:"version"`     // 申请使用的版本号
	Homelicense string `json:"homelicense"`
	Username    string `json:"username"`
}

// ResponseApply 管理员响应用户申请license请求所需要的信息
type ResponseApply struct {
	Updateid int32  `json:"updateid"` // 要对info表中的id响应其请求
	Keyid    int32  `json:"keyid"`    // key表中的id号
	License  string `json:"license"`  // 生成好的license
}

// LicenseInfo 用于管理员查询license_infos全表信息
type LicenseInfo struct {
	ID            int32     `json:"id"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
	StaticID      int32     `json:"static_id"`
	KeyID         int32     `json:"key_id"`
	Quota         uint16    `json:"quota"`
	Version       string    `json:"version"`
	HomeLicense   string    `json:"home_license"`
	License       string    `json:"license"`
	ExpiredTime   time.Time `json:"expired_time"`
	ServicePeriod uint16    `json:"service_period"`
	Status        bool      `json:"status"`
	Username      string    `json:"username"`
}

// UserQueryApplyStatus 查询传入的这个username所申请的状态是否通过
func (l *LicenseInfoModule) UserQueryApplyStatus(username string) ([]ApplyStatus, error) {
	db := dbpool.Pool().DB

	// 记录一个username申请的所有licens状态
	var applystatus []ApplyStatus

	// 找到这个username所在的info表中匹配的所有记录
	var info []models.LicenseInfo
	err := db.Where("username = ?", username).Find(&info).Error
	if err != nil {
		return applystatus, errors.WithMessagef(err,
			"licensedb.UserQueryApplyStatus query apply state and license info into dbtable:license_infos, username = %s", username)
	}

	for _, v := range info {
		var apply ApplyStatus
		apply.Status = v.Status
		apply.License = v.License
		apply.ExpiredTime = v.ExpiredTime
		applystatus = append(applystatus, apply)
	}

	return applystatus, nil

}

// AdminProcessUserApply 管理员根据传入的密钥对ID，version，生成好的License来处理用户申请的License请求
// 并且根据用户申请的服务期限时间生成新的过期时间
func (l *LicenseInfoModule) AdminProcessUserApply(response ResponseApply) error {
	db := dbpool.Pool().DB

	// 查询要更新的记录申请的服务期限
	var info models.LicenseInfo
	db.First(&info, response.Updateid)

	needhour := fmt.Sprintf("%dh", info.ServicePeriod)
	// 使用当前时间加上申请的服务时间得到此License的过期时间
	now := time.Now()
	dd, _ := time.ParseDuration(needhour)
	expiredtime := now.Add(dd)

	// 更新外键, 过期时间，homelicense，license，申请状态，更新时间
	err := db.Exec("update license_infos set key_id = ?, license = ?, update_time = ?, expired_time = ?, status = ? where id = ?",
		response.Keyid, response.License, time.Now(), expiredtime, true, response.Updateid).Error
	if err != nil {
		return errors.WithMessagef(err,
			"licensedb.AdminProcessUserApply update license apply info into dbtable:license_infos, username:%s", info.Username)
	}
	return nil
}

// AdminQueryAllInfo 管理员查询全表信息返回
func (l *LicenseInfoModule) AdminQueryAllInfo() ([]LicenseInfo, error) {
	db := dbpool.Pool().DB

	var info []models.LicenseInfo
	err := db.Find(&info).Error
	if err != nil {
		return nil, errors.WithMessage(err, "licensedb.AdminQueryAllInfo query all info into dbtable:license_infos")
	}

	var respinfo []LicenseInfo
	for _, v := range info {
		var resp LicenseInfo
		resp.CreateTime = v.CreateTime
		resp.ExpiredTime = v.ExpiredTime
		resp.HomeLicense = v.HomeLicense
		resp.ID = v.ID
		resp.KeyID = v.KeyID
		resp.License = v.License
		resp.Quota = v.Quota
		resp.ServicePeriod = v.ServicePeriod
		resp.StaticID = v.StaticID
		resp.UpdateTime = v.UpdateTime
		resp.Username = v.Username
		respinfo = append(respinfo, resp)
	}

	return respinfo, nil
}

// AdminQueryLastInfo 管理员查询最后一条信息返回
func (l *LicenseInfoModule) AdminQueryLastInfo() (LicenseInfo, error) {
	db := dbpool.Pool().DB

	var v models.LicenseInfo
	var respinfo LicenseInfo
	err := db.Find(&v).Error
	if err != nil {
		return respinfo, errors.WithMessage(err, "licensedb.AdminQueryLastInfo query last info into dbtable:license_infos")
	}

	respinfo.CreateTime = v.CreateTime
	respinfo.ExpiredTime = v.ExpiredTime
	respinfo.HomeLicense = v.HomeLicense
	respinfo.ID = v.ID
	respinfo.KeyID = v.KeyID
	respinfo.License = v.License
	respinfo.Quota = v.Quota
	respinfo.ServicePeriod = v.ServicePeriod
	respinfo.StaticID = v.StaticID
	respinfo.UpdateTime = v.UpdateTime
	respinfo.Username = v.Username
	return respinfo, nil
}

// UserInsertInfo 用户插入[license_static.coporation,license_static.uuid],quota,version,service_period, home_license
// 通过此插入方法同步更新license_static表的coporation, uuid信息
func (l *LicenseInfoModule) UserInsertInfo(apply UserApply) error {
	db := dbpool.Pool().DB

	// 执行事务，先在license_static表中插入uuid，coporation等字段，再查找出license_static最后一次insert的id，
	// 将这个id和license_info需要的 quota, version, service_period等字段插入到license_info中
	return db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行对static表的插入
		if err := tx.Create(&models.LicenseStatic{
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
			UUID:        apply.Uuid,
			Corporation: apply.Corporation,
		}).Error; err != nil {
			// 返回任何错误都会回滚事务
			return errors.WithMessagef(err, "licensedb.UserInsertInfo insert uuid:%s, corporation:%s into dbtable:license_statics",
				apply.Uuid, apply.Corporation)
		}

		var result models.LicenseStatic
		// 查询刚刚插入license_static的id == max(id) from license_statics
		tx.Last(&result)

		if err := tx.Create(&models.LicenseInfo{
			StaticID:      result.ID,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			Quota:         apply.Quota,
			ServicePeriod: apply.Period,
			Version:       apply.Version,
			HomeLicense:   apply.Homelicense,
			Username:      apply.Username,
		}).Error; err != nil {
			return errors.WithMessagef(err, "licensedb.UserInsertInfo insert info into dbtable:license_infos, username = %s",
				apply.Username)
		}
		// 返回 nil 提交事务
		return nil
	})
}

// AddLinceseAllInfo 插入license_info的所有字段信息
func (l *LicenseInfoModule) AddLinceseAllInfo(info models.LicenseInfo) error {
	db := dbpool.Pool().DB

	newInfo := models.LicenseInfo{
		CreateTime:    info.CreateTime,
		UpdateTime:    info.UpdateTime,
		ExpiredTime:   info.ExpiredTime,
		License:       info.License,
		HomeLicense:   info.HomeLicense,
		KeyID:         info.KeyID,
		StaticID:      info.StaticID,
		ServicePeriod: info.ServicePeriod,
		Version:       info.Version,
		Quota:         info.Quota,
		Status:        info.Status,
		Username:      info.Username,
	}

	err := db.Create(newInfo).Error
	return errors.WithMessage(err, "licensedb.AddLinceseAllInfo insert full record into dbtable:license_infos")
}
