package models

import "time"

// LicenseInfo license_info
type LicenseInfo struct {
	ID            int32         `gorm:"AUTO_INCREMENT;primary_key;column:id;type:mediumint(8) unsigned;not null" json:"id"`
	CreateTime    time.Time     `gorm:"column:create_time;type:datetime;not null" json:"create_time"`
	UpdateTime    time.Time     `gorm:"column:update_time;type:datetime;not null" json:"update_time"`
	StaticID      int32         `gorm:"index:fk_license_info_license_info_1;column:static_id;type:mediumint(8) unsigned" json:"static_id"`
	LicenseStatic LicenseStatic `gorm:"association_foreignkey:static_id;foreignkey:id" json:"license_static_list"`
	KeyID         int32         `gorm:"index:fk_license_info_license_info_2;column:key_id;type:mediumint(8) unsigned" json:"key_id"`
	LicenseKey    LicenseKey    `gorm:"association_foreignkey:key_id;foreignkey:id" json:"license_key_list"`
	Quota         uint16        `gorm:"column:quota;type:smallint(5) unsigned;not null" json:"quota"`
	Version       string        `gorm:"column:version;type:varchar(50);not null" json:"version"`
	HomeLicense   string        `gorm:"column:home_license;type:varchar(50);not null" json:"home_license"`
	License       string        `gorm:"column:license;type:text;" json:"license"`
	ExpiredTime   time.Time     `gorm:"column:expired_time;type:date" json:"expired_time"`
	ServicePeriod uint16        `gorm:"column:service_period;type:smallint(5) unsigned;not null" json:"service_period"`
	Status        bool          `gorm:"column:status;type:tinyint(1) unsigned;not null;default:0" json:"status"`
	Username      string        `gorm:"column:username;type:varchar(30);not null" json:"username"`
}
