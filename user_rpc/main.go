package main

import (
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"user_rpc/handler"
	"user_rpc/pkg/config"
	"user_rpc/pkg/database"
	"user_rpc/pkg/logger"
	"user_rpc/pkg/redis"
	"user_rpc/proto"
)

func init() {
	// 初始化config logger database redis
	config.SetupConfig()
	logger.SetupLogger()
	database.SetupDatabase()
	redis.SetupRedis()
}

func destroy() {
	// 关闭 database redis
	database.Close()
	redis.Close()
}

func main() {
	// 监听端口
	lis, err := net.Listen("tcp", config.GetString("RPC_SERVER_ADDR"))
	if err != nil {
		logger.Error("main", "listen error", config.GetString("RPC_SERVER_ADDR"))
		panic(err)
	}
	logger.Debug("main", "listen success", config.GetString("RPC_SERVER_ADDR"))

	// 注册并启动服务
	rpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(rpcServer, &handler.UserHandler{})
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			logger.Error("main", "rpc serve err", err)
		}
	}()

	// 监听服务退出状态
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 关闭rpc、数据库、缓存连接
	logger.Info("main", "shutdown server")
	rpcServer.GracefulStop()
	destroy()
	logger.Info("main", "server exit")
}
