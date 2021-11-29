package gprc

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func CloseGRPCConnect(conn *grpc.ClientConn) {
	if conn != nil {
		if err := conn.Close(); err != nil {
			zap.L().Error("close gprc err", zap.Any("error", err))
		}
	}
}

func GetGRPCConnect(ctx context.Context, target string) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, target, grpc.WithChainUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		grpc_zap.UnaryClientInterceptor(zap.L()),
		grpc_opentracing.UnaryClientInterceptor(),
	)),
		grpc.WithInsecure())
	if err != nil {
		zap.L().Error("grpc.DialContext conn error ", zap.String("target", target), zap.Any("error", err))
		return nil
	}
	return conn
}
