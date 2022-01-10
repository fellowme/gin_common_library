package logger

import (
	"context"
	"errors"
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
	switch {
	case err != nil && l.LogConfig.LogLevel >= gorm_logger.Error && (!l.LogConfig.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		l.ZapLogger.Error("trace", zap.Error(err), zap.Duration("cost_time", costTime), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogConfig.SlowThreshold != 0 && costTime > l.LogConfig.SlowThreshold && l.LogConfig.LogLevel >= gorm_logger.Warn:
		sql, rows := fc()
		l.ZapLogger.Warn("trace", zap.Duration("cost_time", costTime), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogConfig.LogLevel >= gorm_logger.Info:
		sql, rows := fc()
		l.ZapLogger.Info("trace", zap.Duration("cost_time", costTime), zap.Int64("rows", rows), zap.String("sql", sql))
	}

}
