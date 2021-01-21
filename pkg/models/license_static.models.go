package models

import "time"

// LicenseStatic license_static
type LicenseStatic struct {
	ID          int32     `gorm:"AUTO_INCREMENT;primary_key;column:id;type:mediumint(8) unsigned;not null" json:"id"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;not null" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;not null" json:"update_time"`
	UUID        string    `gorm:"column:uuid;type:varchar(50);not null" json:"uuid"`
	Corporation string    `gorm:"column:corporation;type:varchar(60);not null" json:"corporation"`
	Extension   string    `gorm:"column:extension;type:varchar(100);default:'null'" json:"extension"`
}
