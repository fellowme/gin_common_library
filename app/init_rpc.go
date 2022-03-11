package app

import (
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func InitRpcServer(configPath, serverName string) *grpc.Server {
	initCommonExtend(configPath, serverName)
	gRpcService := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zap.L()),
		)))
	return gRpcService
}

func CreateRpcServer(gRpcService *grpc.Server) {
	defer gRpcService.GracefulStop()
	defer deferClose()
	endPoint := fmt.Sprintf("%s:%d", gin_config.ServerConfigSettings.Server.ServerHost,
		gin_config.ServerConfigSettings.Server.ServerRpcPort)
	listener, err := net.Listen("tcp", endPoint)
	if err != nil {
		zap.L().Error("rpc listener error", zap.Any("error", err))
		return
	}
	err = gRpcService.Serve(listener)
	zap.L().Info("grpc server starting ...")
	if err != nil {
		zap.L().Error("rpc Serve error", zap.Any("error", err))
		return
	}
}
