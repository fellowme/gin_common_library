package app

import (
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	gin_jaeger "github.com/fellowme/gin_common_library/jaeger"
	gin_logger "github.com/fellowme/gin_common_library/logger"
	gin_mysql "github.com/fellowme/gin_common_library/mysql"
	gin_redis "github.com/fellowme/gin_common_library/redis"
	gin_translator "github.com/fellowme/gin_common_library/translator"
	gin_util "github.com/fellowme/gin_common_library/util"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func initExtend(configPath, serverName string) {
	path := gin_util.GetPath()
	gin_config.InitConfig(path+configPath, serverName)
	gin_logger.InitServerLogger(path)
	gin_logger.InitRecoveryLogger(path)
	gin_jaeger.InitJaegerTracer()
	gin_translator.InitTranslator()
	gin_redis.InitRedis()
	gin_mysql.InitMysqlMap()
}

func CreateRpcServer(configPath, serverName string) {
	initExtend(configPath, serverName)
	gRpcService := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zap.L()),
		)))
	defer gRpcService.GracefulStop()
	defer DeferClose()
	endPoint := fmt.Sprintf("%s:%d", gin_config.ServerConfigSettings.Server.ServerHost,
		gin_config.ServerConfigSettings.Server.ServerRpcPort)
	listener, err := net.Listen("tcp", endPoint)
	if err != nil {
		zap.L().Error("rpc listener error", zap.Any("error", err), zap.String("server_name", serverName))
		return
	}
	err = gRpcService.Serve(listener)
	zap.L().Info("grpc server 启动")
	if err != nil {
		zap.L().Error("rpc Serve error", zap.Any("error", err), zap.String("server_name", serverName))
		return
	}
}

func deferClose() {
	gin_mysql.CloseMysqlConnect()
	gin_jaeger.IoCloser()
	gin_redis.CloseRedisPool()
}
