package middleware

import (
	"lisence/pkg/define"
	"lisence/pkg/handler"
	"lisence/pkg/thirdurl/auth"
	"strings"

	"github.com/kataras/iris/v12"
)

func CheckToken(ctx iris.Context) {
	if ctx.RequestPath(false) == "/login" {
		ctx.Next()
		return
	}

	tokenstr := ctx.Request().Header.Get("Authorization")
	if tokenstr == "" {
		handler.ResponseErrMsg(ctx, define.StatusForbidden, "token非法")
		return
	}
	tokens := strings.Split(tokenstr, " ")
	if len(tokens) != 2 {
		handler.ResponseErrMsg(ctx, define.StatusForbidden, "token非法")
		return
	}

	if u, err := auth.TokenWithSession(tokens[1]); err != nil {
		handler.ResponseErrMsg(ctx, define.StatusForbidden, "token非法")
		return
	} else {
		ctx.Values().Set(define.ReqUserKey, u.Name)
		ctx.Next()
	}
}
