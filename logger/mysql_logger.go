package logger

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	ZapLogger *zap.Logger
	LogConfig gormlogger.Config
}

func NewSqlLogger(zapLogger *zap.Logger, logConfig gormlogger.Config) Logger {
	return Logger{
		ZapLogger: zapLogger,
		LogConfig: logConfig,
	}
}

func (l Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.LogConfig.LogLevel = level
	return l
}
func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogConfig.LogLevel < gormlogger.Info {
		return
	}
	l.ZapLogger.Sugar().Debugf(str, args...)
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogConfig.LogLevel < gormlogger.Warn {
		return
	}
	l.ZapLogger.Sugar().Warnf(str, args...)
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogConfig.LogLevel < gormlogger.Error {
		return
	}
	l.ZapLogger.Sugar().Errorf(str, args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogConfig.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogConfig.LogLevel >= gormlogger.Error && (!l.LogConfig.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		l.ZapLogger.Error("trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogConfig.SlowThreshold != 0 && elapsed > l.LogConfig.SlowThreshold && l.LogConfig.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		l.ZapLogger.Warn("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogConfig.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		l.ZapLogger.Debug("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}
