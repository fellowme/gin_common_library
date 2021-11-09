package logger

import (
	gin_config "github.com/fellowme/gin_commom_library/config"
	gin_util "github.com/fellowme/gin_commom_library/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func initLogger(loggerConfig gin_config.LoggerConfig, serverName string) *zap.Logger {
	filename := serverName + time.Now().Format("20060102") + ".log"
	hook := lumberjack.Logger{
		Filename:   filename,                      // 日志文件路径
		MaxSize:    loggerConfig.LoggerMaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: loggerConfig.LoggerMaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     loggerConfig.LoggerMaxAge,     // 文件最多保存多少天
		Compress:   loggerConfig.LoggerIsCompress, // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    "logger",
		CallerKey:  "line_num",
		MessageKey: "message",
		//StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.Level(loggerConfig.LoggerLevelInt))

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddStacktrace(zapcore.Level(loggerConfig.LoggerLevelInt))
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", serverName))
	// 构造日志
	ZapLogger := zap.New(core, caller, development, filed)
	ZapLogger.Info("zap_log 初始化成功")
	//zap.ReplaceGlobals(ZapLogger)
	return ZapLogger
}

func RecoveryWithZap(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.String("time", gin_util.NowTimeToString()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.String("time", gin_util.NowTimeToString()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
