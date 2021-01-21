package middleware

import (
	// "github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/jwt"
)

/**
 * 验证 jwt
 * @method JwtHandler
 */
func JwtHandler() *jwt.Middleware {
	// 验证 jwt 的 token 的方法
	return jwt.New(jwt.Config{
		// 从请求头的Authorization字段中提取，这个是默认值
		Extractor: jwt.FromAuthHeader, // Extractor: jwt.FromParameter("token"),

		// 设置一个函数返回秘钥，关键在于return []byte("这里设置秘钥")
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// 自己加密的密钥/盐
			return []byte("TroilaSecret"), nil
		},

		// ContextKey 请求中用户(&令牌)信息所在的属性的名称，如果是 JWT 将会被存储,默认值为 jwt
		// 设置之后，middelware将验证是否使用特定的签名算法对令牌进行了签名，如果签名方法不是常量，可以使用ValidationKeyGetter回调来实现其他检查，默认值 nil
		SigningMethod: jwt.SigningMethodHS256,
		// 验证 令牌是否过期，如果过期返回 过期错误，默认值 false
		Expiration: true,
	})
}
