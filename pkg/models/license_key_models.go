package models

import "time"

// LicenseKey license_key
type LicenseKey struct {
	ID          int32     `gorm:"AUTO_INCREMENT;primary_key;column:id;type:mediumint(8) unsigned;not null" json:"id"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;not null" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;not null" json:"update_time"`
	PublicKey   string    `gorm:"column:public_key;type:text;not null" json:"public_key"`
	PrivateKey  string    `gorm:"column:private_key;type:text;not null" json:"private_key"`
	Username    string    `gorm:"column:username;type:varchar(30);not null" json:"username"`
	Description string    `gorm:"column:description;type:varchar(100);default:'null'" json:"description"`
}
