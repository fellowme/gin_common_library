package app

import (
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	grpc_consul "github.com/fellowme/gin_common_library/consul"
	gin_middleware "github.com/fellowme/gin_common_library/middleware"
	"github.com/fellowme/gin_common_library/service_health"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func InitRpcServer(configPath, serverName string) *grpc.Server {
	initCommonExtend(configPath, serverName)
	grpc_consul.InitConsulClient()
	gRpcService := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zap.L()),
			gin_middleware.UnaryServerInterceptor(),
		)))
	return gRpcService
}

func CreateRpcServer(gRpcService *grpc.Server, consul grpc_consul.ServiceConsul) {
	defer func() {
		gRpcService.GracefulStop()
		deferClose()
		grpc_consul.UnRegisterConsul(consul.Id)
	}()
	errChan := make(chan error)
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		grpc_health_v1.RegisterHealthServer(gRpcService, &service_health.HealthService{})
		endPoint := fmt.Sprintf("%s:%d", gin_config.ServerConfigSettings.Server.ServerHost,
			gin_config.ServerConfigSettings.Server.ServerRpcPort)
		listener, err := net.Listen("tcp", endPoint)
		if err != nil {
			errChan <- err
			zap.L().Error("rpc listener error", zap.Any("error", err))
			return
		}
		err = grpc_consul.RegisterGrpcConsul(consul)
		if err != nil {
			errChan <- err
			zap.L().Error("rpc RegisterConsul error", zap.Any("error", err))
			return
		}
		zap.L().Info("grpc server starting ...")
		err = gRpcService.Serve(listener)
		if err != nil {
			errChan <- err
			zap.L().Error("rpc start error", zap.Any("error", err))
		}
	}()
	select {
	case err := <-errChan:
		zap.L().Error("rpc Serve error", zap.Any("error", err))
	case data := <-stopChan:
		zap.L().Info("rpc stop", zap.Any("signal", data))
	}
}
