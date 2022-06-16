package logger

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

type Logger struct {
	ZapLogger *zap.Logger
	LogConfig gorm_logger.Config
}

func NewSqlLogger(zapLogger *zap.Logger, logConfig gorm_logger.Config) Logger {
	return Logger{
		ZapLogger: zapLogger,
		LogConfig: logConfig,
	}
}

func (l Logger) LogMode(level gorm_logger.LogLevel) gorm_logger.Interface {
	l.LogConfig.LogLevel = level
	return l
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogConfig.LogLevel >= gorm_logger.Info {
		l.ZapLogger.Sugar().Infof(str, args...)
	}
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogConfig.LogLevel >= gorm_logger.Warn {
		l.ZapLogger.Sugar().Warnf(str, args...)
	}

}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogConfig.LogLevel >= gorm_logger.Error {
		l.ZapLogger.Sugar().Errorf(str, args...)
	}

}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogConfig.LogLevel <= 0 {
		return
	}
	costTime := time.Since(begin)
	sql, rows := fc()
	switch {
	case err != nil && l.LogConfig.LogLevel >= gorm_logger.Error && (!l.LogConfig.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.ZapLogger.Error("trace", zap.Error(err), zap.Float64("cost_time", costTime.Seconds()), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogConfig.SlowThreshold != 0 && costTime > l.LogConfig.SlowThreshold && l.LogConfig.LogLevel >= gorm_logger.Warn:
		l.ZapLogger.Warn("trace", zap.Float64("cost_time", costTime.Seconds()), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogConfig.LogLevel >= gorm_logger.Info:
		l.ZapLogger.Info("trace", zap.Float64("cost_time", costTime.Seconds()), zap.Int64("rows", rows), zap.String("sql", sql))
	}
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		span := opentracing.StartSpan(
			"gorm_action_trace",
			opentracing.ChildOf(parentSpan.Context()),
			opentracing.Tags{
				"cost_time": costTime.Seconds(),
				"rows":      rows,
				"sql":       sql,
			},
			ext.SpanKindProducer,
		)
		defer span.Finish()
	}

}
