package models

import "time"

// CmsUser user
type CmsUser struct {
	ID         int32     `gorm:"AUTO_INCREMENT;primary_key;column:id;type:mediumint(8) unsigned;not null" json:"id"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;not null" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime;not null" json:"update_time"`
	Username   string    `gorm:"column:username;type:varchar(30);not null" json:"username"`
	Password   string    `gorm:"column:password;type:varchar(128);not null" json:"password"`
	IsAdmin    bool      `gorm:"column:is_admin;type:tinyint(1) unsigned;not null;default:0" json:"is_admin"`
}
