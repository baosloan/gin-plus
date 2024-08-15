package cmd

import (
	"context"
	"errors"
	"gin-plus/api/routes"
	"gin-plus/config"
	"gin-plus/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Run: func(cmd *cobra.Command, args []string) {
		//配置文件初始化
		config.Init(configPath)

		//初始化日志
		log.Init()
		defer func() {
			_ = log.Sync()
		}()

		//初始化路由
		gin.SetMode(config.Conf.Mode)
		router := routes.Init()

		//启动服务
		server := http.Server{
			Addr:           ":" + strconv.Itoa(config.Conf.Server.Port),
			Handler:        router,
			ReadTimeout:    config.Conf.Server.ReadTimeout * time.Second,
			WriteTimeout:   config.Conf.Server.WriteTimeout * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		log.Info("启动http服务器")
		go func() {
			//开启一个goroutine启动服务
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("http服务启动失败", zap.Error(err))
			}
		}()
		// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
		quit := make(chan os.Signal, 1)
		// kill 默认会发送 syscall.SIGTERM 信号
		// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
		// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
		// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit //阻塞在此，当接收到上述两种信号时才会往下执行
		log.Info("Shutdown Server ...")
		//创建一个5秒超时的context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown failed", zap.Error(err))
		}
		log.Info("Server exited")
	},
}

var configPath string

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&configPath, "config", "c", "", "配置文件路径")
}
