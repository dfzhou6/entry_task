package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
	"user_web/client"
	"user_web/pkg/config"
	"user_web/pkg/logger"
	"user_web/pkg/redis"
	"user_web/route"
)

func init() {
	// 初始化logger、config、redis、client
	config.SetupConfig()
	logger.SetupLogger()
	redis.SetupRedis()
	client.SetupClient()
}

func destroy() {
	// 关闭redis、client连接
	redis.Close()
	client.Close()
}

func main() {
	// 注册并绑定http服务
	r := gin.New()
	route.RegisterHandler(r)
	gin.SetMode(gin.ReleaseMode)
	srv := &http.Server{
		Addr:    config.GetString("WEB_SERVER_ADDR"),
		Handler: r,
	}

	// 启动http服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("main", "listenAndServe err", err)
		}
	}()

	// 监听退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 关闭服务和连接
	logger.Info("main", "shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("main", "shutdown server err", err)
	}
	destroy()
	logger.Info("main", "server exit")
}
