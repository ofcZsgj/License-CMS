package handler

import (
	"lisence/pkg/dao/licensedb"
	"lisence/pkg/define"
	"lisence/pkg/service"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"

	log "github.com/sirupsen/logrus"
)

// QueryKeyInfo 管理员可以查询所有的公私钥对信息
func QueryKeyInfo(ctx iris.Context) {
	jwtInfo := ctx.Values().Get("jwt").(*jwt.Token)
	isadmin := jwtInfo.Claims.(jwt.MapClaims)["isadmin"].(bool)
	if !isadmin {
		ResponseErrMsg(ctx, define.StatusForbidden, "is not admin, forbidden token")
		return
	}

	retdata, err := service.GetSingleTon().KeyInfo()
	if err != nil {
		ResponseErr(ctx, err)
		return
	}
	ResponseOk(ctx, retdata)
	log.Info("query key reponse ok")
}

// InsertKeyInfo 管理员插入一条新的key表数据
func InsertKeyInfo(ctx iris.Context) {
	jwtInfo := ctx.Values().Get("jwt").(*jwt.Token)
	isadmin := jwtInfo.Claims.(jwt.MapClaims)["isadmin"].(bool)
	if !isadmin {
		ResponseErrMsg(ctx, define.StatusForbidden, "is not admin, forbidden token")
		return
	}

	k := licensedb.KeyPairInfo{}
	err := ctx.ReadJSON(&k)
	if err != nil {
		ResponseErrMsg(ctx, define.StatusBadRequest, "json format error")
		return
	}

	if k.Description == "" || k.Username == "" {
		ResponseErrMsg(ctx, define.StatusBadRequest, "exist nil filed")
		return
	}

	err = service.GetSingleTon().InsertKeyInfomation(k)
	if err != nil {
		ResponseErr(ctx, err)
		return
	}
	ResponseOk(ctx, "OK")
	log.Infof("Admin insert key, username:%s insert key reponse ok", k.Username)
}

// UserApplyLicense 用户申请license
func UserApplyLicense(ctx iris.Context) {
	applyinfo := licensedb.UserApply{}
	err := ctx.ReadJSON(&applyinfo)
	if err != nil {
		ResponseErrMsg(ctx, define.StatusBadRequest, "json format error")
		return
	}

	if applyinfo.Corporation == "" || applyinfo.Homelicense == "" || applyinfo.Period == 0 ||
		applyinfo.Quota == 0 || applyinfo.Username == "" || applyinfo.Uuid == "" || applyinfo.Version == "" {
		ResponseErrMsg(ctx, define.StatusBadRequest, "exist nil filed")
		return
	}

	err = service.GetSingleTon().ApplyLicense(applyinfo)
	if err != nil {
		ResponseErr(ctx, err)
		return
	}
	ResponseOk(ctx, "OK")
	log.Infof("user apply license, username:%s reponse ok", applyinfo.Username)
}

// UserQueryLicense 用户查询license
func UserQueryLicense(ctx iris.Context) {
	username := ctx.Params().GetString("username")
	if username == "" {
		ResponseErrMsg(ctx, define.StatusBadRequest, "params is nil")
		return
	}

	status, err := service.GetSingleTon().QueryLicense(username)
	if err != nil {
		ResponseErr(ctx, err)
		return
	}
	ResponseOk(ctx, status)
	log.Infof("user query license, username:%s query license reponse ok", username)
}

// AdminQueryLicense 管理员查询所有的申请
func AdminQueryLicense(ctx iris.Context) {
	jwtInfo := ctx.Values().Get("jwt").(*jwt.Token)
	isadmin := jwtInfo.Claims.(jwt.MapClaims)["isadmin"].(bool)
	if !isadmin {
		ResponseErrMsg(ctx, define.StatusForbidden, "is not admin, forbidden token")
		return
	}

	apply, err := service.GetSingleTon().QueryLicenseApply()
	if err != nil {
		ResponseErr(ctx, err)
		return
	}
	ResponseOk(ctx, apply)
	log.Infof("admin query license apply info reponse ok")
}

// AdminProcessApply 管理员处理一条指定的license申请
func AdminProcessApply(ctx iris.Context) {
	jwtInfo := ctx.Values().Get("jwt").(*jwt.Token)
	isadmin := jwtInfo.Claims.(jwt.MapClaims)["isadmin"].(bool)
	if !isadmin {
		ResponseErrMsg(ctx, define.StatusForbidden, "is not admin, forbidden token")
		return
	}

	respapply := licensedb.ResponseApply{}
	err := ctx.ReadJSON(&respapply)
	if err != nil {
		ResponseErrMsg(ctx, define.StatusBadRequest, "json format error")
		return
	}

	if respapply.Keyid == 0 || respapply.Updateid == 0 {
		ResponseErrMsg(ctx, define.StatusBadRequest, "exist nil filed")
		return
	}

	err = service.GetSingleTon().ProcessLicenseApply(respapply)
	if err != nil {
		ResponseErr(ctx, err)
		return
	}
	ResponseOk(ctx, "OK")
	log.Infof("admin process license apply, updateid:%d, keyid:%d reponse ok", respapply.Updateid, respapply.Keyid)
}

// Regist 注册
func Regist(ctx iris.Context) {
	regapply := licensedb.UserInfo{}
	err := ctx.ReadJSON(&regapply)
	if err != nil {
		ResponseErrMsg(ctx, define.StatusBadRequest, "json format error")
		return
	}

	if regapply.Username == "" || regapply.Password == "" {
		ResponseErrMsg(ctx, define.StatusBadRequest, "exist nil filed")
		return
	}

	err = service.GetSingleTon().Registration(regapply)
	if err != nil {
		ResponseErr(ctx, err)
		return
	}

	ResponseOk(ctx, "OK")
	log.Infof("regist success, username:%s, reponse ok", regapply.Username)
}

// Login 登录
func Login(ctx iris.Context) {
	logininfo := licensedb.UserInfo{}
	err := ctx.ReadJSON(&logininfo)
	if err != nil {
		ResponseErrMsg(ctx, define.StatusBadRequest, "json format error")
		return
	}

	if logininfo.Username == "" || logininfo.Password == "" {
		ResponseErrMsg(ctx, define.StatusBadRequest, "exist nil filed")
		return
	}

	flag, err := service.GetSingleTon().LoginCMS(logininfo)
	if err != nil {
		ResponseErr(ctx, err)
		return
	}

	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// 用户名
		"username": logininfo.Username,
		// 是否是管理员
		"isadmin": flag,
		// 签发时间
		"iat": time.Now().Unix(),
		// 添加过期时间 24H
		"exp": time.Now().Add(1 * time.Hour * time.Duration(24)).Unix(),
	})
	tokenStr, err := token.SignedString([]byte("TroilaSecret"))
	if err != nil {
		log.Info(err)
		ResponseErrMsg(ctx, define.StatusBadRequest, "generate jwt token fail")
		return
	}

	ResponseOk(ctx, tokenStr)
	log.Infof("login success, username:%s, isadmin:%t", logininfo.Username, flag)
}
