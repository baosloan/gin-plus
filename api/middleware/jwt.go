package middleware

import (
	"errors"
	pkgJWT "gin-plus/pkg/jwt"
	"gin-plus/pkg/resp"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		//客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URL中
		//优先从请求头中获取Token
		token := c.GetHeader("token")
		if token == "" {
			//token中获取不到，则尝试从请求体中获取Token
			token = c.DefaultPostForm("token", "")
			if token == "" {
				//请求头和请求体中都获取不到，从URL获取
				token = c.DefaultQuery("token", "")
				if token == "" {
					c.JSON(http.StatusUnauthorized, gin.H{
						"code": http.StatusUnauthorized,
						"msg":  http.StatusText(http.StatusUnauthorized),
					})
					c.Abort()
					return
				}
			}
		}
		claims, err := pkgJWT.Parse(token)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				resp.Fail(c, resp.CodeTokenExpired)
			case errors.Is(err, jwt.ErrTokenMalformed):
				resp.Fail(c, resp.CodeTokenMalformed)
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				resp.Fail(c, resp.CodeTokenSignatureInvalid)
			case errors.Is(err, jwt.ErrTokenInvalidAudience):
				resp.Fail(c, resp.CodeTokenInvalidAudience)
			case errors.Is(err, jwt.ErrTokenUsedBeforeIssued):
				resp.Fail(c, resp.CodeTokenUsedBeforeIssued)
			case errors.Is(err, jwt.ErrTokenInvalidIssuer):
				resp.Fail(c, resp.CodeTokenInvalidIssuer)
			case errors.Is(err, jwt.ErrTokenInvalidSubject):
				resp.Fail(c, resp.CodeTokenInvalidSubject)
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				resp.Fail(c, resp.CodeTokenNotValidYet)
			case errors.Is(err, jwt.ErrTokenInvalidId):
				resp.Fail(c, resp.CodeTokenInvalidId)
			case errors.Is(err, jwt.ErrTokenInvalidClaims):
				resp.Fail(c, resp.CodeTokenInvalidClaims)
			default:
				resp.Fail(c, resp.CodeUnknown)
			}
			c.Abort()
			return
		}
		//将当前请求的用户ID保存到请求的上下文中，后续的处理函数可以通过c.Get('user_id')来获取当前请求的用户信息
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
