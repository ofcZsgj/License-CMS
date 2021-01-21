package router

import (
	"lisence/pkg/handler"
	"lisence/pkg/middleware"
	"net/http"
)

////////////////////////////////
// router
func (a *API) InitRouter() *API {
	//init middleware here
	a.OnErrorCodeLog(http.StatusNotFound, middleware.Log404Error)
	a.OnErrorCodeLog(http.StatusInternalServerError, middleware.Log500Error)
	a.OnErrorCodeLog(http.StatusBadGateway, middleware.Log502Error)
	// 保持连接
	a.SetMiddleware(middleware.IAmAlive)
	// access log, 对请求和应答都做日志记录
	a.SetMiddleware(handler.RequestLog)

	// check session
	//a.SetMiddleware(middleware.CheckToken)

	// keepalived api
	a.Get("/do_not_delete.html", middleware.IAmAlive)
	//a.Any("/-/reload", handler.Reload)

	//global api demo
	{
		// 注册
		a.Post("/regist", handler.Regist)
		// 登录
		a.Post("/login", handler.Login)
	}

	// grouop routing api
	{
		// 业务功能路由
		ga := a.Group("/licenseguard/api/v1", middleware.JwtHandler().Serve)
		// user api路由组
		u := ga.Group("/user")
		// admin api路由组
		r := ga.Group("/admin")

		// 用户申请license
		u.Post("/license", handler.UserApplyLicense)
		// 用户查询申请的license是否审批，license，license过期时间
		u.Get("/license/{username}", handler.UserQueryLicense)

		// 查询所有的Key信息
		r.Get("/key", handler.QueryKeyInfo)
		// 插入新的Key信息
		r.Post("/key", handler.InsertKeyInfo)
		// 管理员查询申请license的请求
		r.Get("/license", handler.AdminQueryLicense)
		// 管理员批准申请license的请求
		r.Post("/license", handler.AdminProcessApply)
	}

	return a
}
