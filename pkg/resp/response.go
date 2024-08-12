package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": CodeSuccess,
		"msg":  codeText[CodeSuccess],
		"data": data,
	})
}

func Fail(c *gin.Context, code int) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  codeText[code],
	})
}

const (
	CodeSuccess     = 0
	CodeUnknown     = 1
	CodeServerError = 500
	CodeNotFound    = 404

	CodeTokenExpired          = 1000
	CodeTokenMalformed        = 1001
	CodeTokenSignatureInvalid = 1002
	CodeTokenInvalidAudience  = 1003
	CodeTokenUsedBeforeIssued = 1004
	CodeTokenInvalidIssuer    = 1005
	CodeTokenInvalidSubject   = 1006
	CodeTokenNotValidYet      = 1007
	CodeTokenInvalidId        = 1008
	CodeTokenInvalidClaims    = 1009
	CodeInvalidParam          = 1000_000
)

var codeText = map[int]string{
	CodeSuccess:               "success",
	CodeUnknown:               "未知错误",
	CodeNotFound:              "404 not found",
	CodeServerError:           "服务器错误",
	CodeTokenExpired:          "令牌已过期",
	CodeTokenMalformed:        "令牌格式不正确",
	CodeTokenSignatureInvalid: "令牌签名不合法",
	CodeTokenInvalidAudience:  "令牌受众无效",
	CodeTokenUsedBeforeIssued: "发行前使用的令牌",
	CodeTokenInvalidIssuer:    "令牌颁发者无效",
	CodeTokenInvalidSubject:   "令牌主题无效",
	CodeTokenNotValidYet:      "令牌尚未生效",
	CodeTokenInvalidId:        "令牌ID无效",
	CodeTokenInvalidClaims:    "令牌声明无效",

	CodeInvalidParam: "参数错误",
}
