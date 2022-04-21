package client

import (
	"google.golang.org/grpc"
	"user_web/pkg/config"
	"user_web/pkg/logger"
	"user_web/proto"
)

var RpcClient proto.UserServiceClient
var internalConn *grpc.ClientConn

// SetupClient 初始化rpcClient
func SetupClient() {
	var err error
	internalConn, err = grpc.Dial(config.GetString("RPC_SERVER_ADDR"), grpc.WithInsecure())
	if err != nil {
		logger.Error("client", "conn client error", err)
		panic(err)
	}
	RpcClient = proto.NewUserServiceClient(internalConn)
	logger.Debug("client", "conn success")
}

// Close 关闭rpcClient连接
func Close() {
	if internalConn == nil {
		return
	}
	if err := internalConn.Close(); err != nil {
		logger.Error("client", "close rpc client error", err)
		panic(err)
	}
	logger.Debug("client", "close success")
}
