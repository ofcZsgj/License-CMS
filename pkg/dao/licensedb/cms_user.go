package licensedb

import (
	"lisence/pkg/dao/dbpool"
	"lisence/pkg/libs/encrypt"
	"lisence/pkg/models"
	"time"

	"github.com/pkg/errors"
)

// UserInfo 用户登录或者注册的信息
type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CmsUserModule cms_users表的数据对象CmsUser操作类
type CmsUserModule struct{}

// UserReg 普通用户注册
func (c *CmsUserModule) UserReg(regapply UserInfo) error {
	db := dbpool.Pool().DB

	// 禁止注册相同的用户名
	if db.Where("username = ?", regapply.Username).First(&models.CmsUser{}).RecordNotFound() == false {
		return errors.Errorf("username:%s is exist", regapply.Username)
	}

	// 加密密码
	lastpwd := encrypt.EncryptPwd(regapply.Username, regapply.Password)

	newUser := models.CmsUser{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Username:   regapply.Username,
		Password:   string(lastpwd),
		IsAdmin:    false,
	}
	err := db.Create(&newUser).Error
	if err != nil {
		return errors.WithMessagef(err, "licensedb.UserReg insert user info into dbtable:cms_users, username = %s", regapply.Username)
	}
	return nil
}

// QueryUserInfo 验证用户名密码并判断是否是管理员
func (c *CmsUserModule) QueryUserInfo(logininfo UserInfo) (bool, error) {
	db := dbpool.Pool().DB
	var user models.CmsUser

	// 生成加密后的密码来和数据库中存储的密码进行比对
	lastpwd := encrypt.EncryptPwd(logininfo.Username, logininfo.Password)
	if err := db.Raw("select * from cms_users where username = ?",
		logininfo.Username).Scan(&user).Error; err != nil {
		// 用户不存在
		return false, errors.WithMessagef(err, "licensedb.QueryUserInfo query username, password into dbtable:cms_users, username = %s", logininfo.Username)
	}

	if user.Password != lastpwd {
		return false, errors.Errorf("password of username:%s is wrong!", user.Username)
	}

	// 判断是管理员/普通用户
	if user.IsAdmin == true {
		return true, nil // 管理员
	} else {
		return false, nil // 普通用户
	}
}

// AddUser 注册管理员/普通用户
func (c *CmsUserModule) AddUser(admin bool, username, password string) error {
	db := dbpool.Pool().DB

	newUser := models.CmsUser{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Username:   username,
		Password:   password,
		IsAdmin:    admin,
	}
	return db.Create(&newUser).Error
}

// CreateCmsUser 为模型`CmsUser`创建表
func (c *CmsUserModule) CreateCmsUser() {
	db := dbpool.Pool().DB

	db.CreateTable(&models.CmsUser{})
}
