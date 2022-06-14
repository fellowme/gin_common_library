package middleware

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/*
	UnaryServerInterceptor 超时 取消
*/
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if ctx.Err() == context.Canceled {
			zap.L().Error("time out Canceled")
			return nil, errors.New("grpc time out")
		}
		resp, err := handler(ctx, req)
		return resp, err
	}
}
