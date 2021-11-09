package gin_remote_service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type Mock func() (body []byte)

type Option func(*option)

type option struct {
	ttl           time.Duration     //超时时间
	header        map[string]string // header 设置
	retryTimes    int               // 重试次数
	retryDelay    time.Duration     // 重试延迟时间
	mock          Mock              // mock值  测试使用
	logger        *zap.Logger       // 远程调用日志
	jaegerContext *gin.Context
}

func getOption() *option {
	return &option{
		ttl:           defaultTTl,
		header:        make(map[string]string),
		retryTimes:    defaultRetryTimes,
		retryDelay:    defaultRetryDelay,
		mock:          nil,
		logger:        zap.L(),
		jaegerContext: nil,
	}
}
func WithTTL(ttl time.Duration) Option {
	return func(o *option) {
		o.ttl = ttl
	}
}

func WithHeader(key, value string) Option {
	return func(opt *option) {
		opt.header[key] = value
	}
}

func WithOnFailedRetry(retryTimes int, retryDelay time.Duration, retryVerify RetryVerify) Option {
	return func(opt *option) {
		opt.retryTimes = retryTimes
		opt.retryDelay = retryDelay
	}
}

func WithMock(m Mock) Option {
	return func(opt *option) {
		opt.mock = m
	}
}

func WithLogger(l *zap.Logger) Option {
	return func(o *option) {
		o.logger = l
	}
}

func WithJaegerContext(jaegerContext *gin.Context) Option {
	return func(o *option) {
		o.jaegerContext = jaegerContext
	}
}
