package jwt

import (
	"gin-plus/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Claims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 如果我们需要记录额外的字段信息，需要自定义结构体
type Claims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

// Sign 签名生成token
func Sign(userId uint, expiredTime time.Duration) (string, error) {
	claims := Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gin-plus",
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Second * expiredTime)},
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//使用指定的secret签名并获得完整编码后的字符串token
	return token.SignedString([]byte(config.Conf.JWT.Secret))
}

// Parse 解析并校验token
func Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(config.Conf.JWT.Secret), nil
	})
	if token != nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, err
}
