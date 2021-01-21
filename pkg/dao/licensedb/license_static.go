package licensedb

import (
	"lisence/pkg/dao/dbpool"
	"lisence/pkg/models"
	"time"
)

// LicenseStaticModule license_statics表的数据对象LicenseStatic操作类
type LicenseStaticModule struct{}

// AddStaticInfo 向license_static表添加数据
func (l *LicenseStaticModule) AddStaticInfo(uuid, corporation, extension string) error {
	db := dbpool.Pool().DB

	newStatic := models.LicenseStatic{
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		UUID:        uuid,
		Corporation: corporation,
		Extension:   extension,
	}
	return db.Create(&newStatic).Error
}

// CreateLicenseStatic 为模型LicenseStatic`创建表
func (l *LicenseStaticModule) CreateLicenseStatic() {
	db := dbpool.Pool().DB
	// 全局禁用表名复数
	//db.SingularTable(true)
	db.CreateTable(&models.LicenseStatic{})
}
