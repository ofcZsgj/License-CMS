package licensedb

import (
	"lisence/pkg/dao/dbpool"
	"lisence/pkg/models"
	"time"

	"github.com/pkg/errors"
)

// LicenseKeyModule license_keys表的数据对象LicenseKey操作类
type LicenseKeyModule struct{}

// KeyPair 密钥对信息
type KeyPairInfo struct {
	ID          int32  `json:"id"`          // 密钥对编号
	Public      string `json:"public"`      // 公钥
	Private     string `json:"private"`     // 私钥
	Username    string `json:"username"`    //这对公私钥的用户名
	Description string `json:"description"` // 描述
}

// KeyPair 密钥对
type KeyPair struct {
	Public  string `json:"public"`  // 公钥
	Private string `json:"private"` // 私钥
}

// AdminQueryAllKey 管理员查询所有的公私钥对
func (l *LicenseKeyModule) AdminQueryAllKey() ([]KeyPairInfo, error) {
	db := dbpool.Pool().DB

	// 找到所有的公私钥，存到key中
	var key []models.LicenseKey
	err := db.Select("id, public_key, private_key, username, description").Find(&key).Error

	if err != nil {
		return nil, errors.WithMessage(err, "licensedb.AdminQueryAllKey query all key info into dbtable:license_keys")
	}

	// 从key中将公私钥字段逐一读入到pair中，再将pair一个一个append到keypair中返回
	var keypair []KeyPairInfo
	for _, v := range key {
		var pair KeyPairInfo
		pair.ID = v.ID
		pair.Public = v.PublicKey
		pair.Private = v.PrivateKey
		pair.Username = v.Username
		pair.Description = v.Description
		keypair = append(keypair, pair)
	}
	return keypair, nil
}

// AdminInsertKey 管理员插入一条新记录,一组公私钥对及这组公私钥对所对应的username和描述
func (l *LicenseKeyModule) AdminInsertKey(pair KeyPairInfo) error {
	db := dbpool.Pool().DB
	var licensekey models.LicenseKey

	licensekey.CreateTime = time.Now()
	licensekey.UpdateTime = time.Now()
	licensekey.PrivateKey = pair.Private
	licensekey.PublicKey = pair.Public
	licensekey.Username = pair.Username
	licensekey.Description = pair.Description

	err := db.Create(&licensekey).Error
	if err != nil {
		return errors.WithMessagef(err, "licensedb.AdminQueryAllKey query all key info into dbtable:license_keys, username = %s", pair.Username)
	}
	return nil
}
