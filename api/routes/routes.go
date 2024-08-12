package routes

import (
	"gin-plus/api/middleware"
	"gin-plus/pkg/resp"
	"github.com/gin-gonic/gin"
	"time"
)

func Init() *gin.Engine {
	//创建一个默认的gin.Engine，使用了Logger()和Recovery()中间件
	r := gin.Default()
	//引入JWT中间件
	r.Use(middleware.JWT())

	r.GET("ping", func(c *gin.Context) {
		time.Sleep(time.Second * 10)
		resp.Success(c, "pong")
	})
	r.NoRoute(func(c *gin.Context) {
		resp.Fail(c, resp.CodeNotFound)
	})
	return r
}
