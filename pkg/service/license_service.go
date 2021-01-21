package service

import (
	"lisence/pkg/dao/licensedb"
	"lisence/pkg/define"
	"lisence/pkg/libs/license"
	"lisence/pkg/libs/rsakey"
	"sync"

	"github.com/sirupsen/logrus"
)

//  singleton service
var (
	once sync.Once
	ser  *LicenseService
)

type LicenseService struct {
}

func GetSingleTon() *LicenseService {
	once.Do(
		func() {
			ser = &LicenseService{}
		})

	return ser
}

// KeyInfo 调用dao/license中的AdminQueryAllKey查询数据库返回key信息
func (d *LicenseService) KeyInfo() ([]licensedb.KeyPairInfo, error) {
	var keymodule licensedb.LicenseKeyModule
	info, err := keymodule.AdminQueryAllKey()
	if err != nil {
		logrus.Errorf("KeyInfo query key info fail, err:%s", err)
		return nil, define.ErrInnerServer("query table license_keys error", err)
	}

	return info, nil
}

// InsertKeyInfomation 调用dao/license中的AdminInsertKey插入一条key信息到数据库
func (d *LicenseService) InsertKeyInfomation(pair licensedb.KeyPairInfo) error {
	// 生成公私钥
	privkey, pubkey, err := rsakey.GenerateRSAKeys()
	if err != nil {
		logrus.Errorf("rsakey.GenerateRSAKeys fail, err:%s", err)
		return define.ErrInnerServer("generate key error", err)
	}
	pair.Public = pubkey
	pair.Private = privkey

	var keymodule licensedb.LicenseKeyModule
	err = keymodule.AdminInsertKey(pair)
	if err != nil {
		logrus.Errorf("InsertKeyInfomation insert key to license_keys fail, err:%s", err)
		return define.ErrInnerServer("insert table license_keys error", err)
	}

	return nil
}

// ApplyLicense 调用dao/license中的UserInsertInfo插入申请license数据，更新info表和static表
func (d *LicenseService) ApplyLicense(applyinfo licensedb.UserApply) error {
	var infomodule licensedb.LicenseInfoModule
	err := infomodule.UserInsertInfo(applyinfo)
	if err != nil {
		logrus.Errorf("ApplyLicense insert license apply fail err:%s", err)
		return define.ErrInnerServer("insert table license_infos, license_statics error", err)
	}

	return nil
}

// ApplyLicense 调用dao/license中的UserQueryApplyStatus查询uername的所有申请记录和license记录
func (d *LicenseService) QueryLicense(username string) ([]licensedb.ApplyStatus, error) {
	var infomodule licensedb.LicenseInfoModule
	status, err := infomodule.UserQueryApplyStatus(username)
	if err != nil {
		logrus.Errorf("UserQueryApplyStatus query license apply status fail, err:%s", err)
		return nil, define.ErrInnerServer("query table license_infos error", err)
	}

	return status, nil
}

// QueryLicenseApply 调用dao/license中的AdminQueryAllInfo查询license_infos表所有的申请记录
func (d *LicenseService) QueryLicenseApply() ([]licensedb.LicenseInfo, error) {
	var infomodule licensedb.LicenseInfoModule
	apply, err := infomodule.AdminQueryAllInfo()
	if err != nil {
		logrus.Errorf("AdminQueryAllInfo query license apply info fail, err:%s", err)
		return nil, define.ErrInnerServer("query table license_infos error", err)
	}

	return apply, nil
}

// ProcessLicenseApply 调用dao/license中的AdminProcessUserApply处理指定记录的license申请
func (d *LicenseService) ProcessLicenseApply(response licensedb.ResponseApply) error {
	// 获得随机生成的license
	response.License = license.RandStringBytes()

	var infomodule licensedb.LicenseInfoModule
	err := infomodule.AdminProcessUserApply(response)
	if err != nil {
		logrus.Errorf("AdminQueryAllInfo query license apply info fail, err:%s", err)
		return define.ErrInnerServer("update table license_infos error", err)
	}

	return nil
}

// Registration 用户注册
func (d *LicenseService) Registration(regapply licensedb.UserInfo) error {
	var usermodule licensedb.CmsUserModule
	err := usermodule.UserReg(regapply)
	if err != nil {
		logrus.Errorf("UserReg insert new user regist info fail, err:%s", err)
		return define.ErrInnerServer("insert table cms_users error", err)
	}

	return nil
}

// LoginCMS 用户登录
func (d *LicenseService) LoginCMS(logininfo licensedb.UserInfo) (bool, error) {
	var usermodule licensedb.CmsUserModule
	flag, err := usermodule.QueryUserInfo(logininfo)
	if err != nil {
		logrus.Errorf("QueryUserInfo query user's username or password fail, err:%s", err)
		return flag, define.ErrInnerServer("query table cms_users error, username or password is not exist", err)
	}

	return flag, nil
}
