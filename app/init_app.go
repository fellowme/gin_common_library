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
	"github.com/fvbock/endless"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

/*
	初始化配置文件
*/
func initCommonExtend(configPath string, serverName string) {
	basePath := gin_util.GetPath()
	gin_config.InitConfig(basePath+configPath, serverName)
	gin_logger.InitServerLogger(basePath)
	gin_logger.InitRecoveryLogger(basePath)
	gin_translator.InitTranslator()
	gin_jaeger.InitJaegerTracer()
	gin_mysql.InitMysqlMap()
	gin_redis.InitRedis()
}

/*
	creatApp 初始化app
*/
func creatApp(configPath, serverName string) *gin.Engine {
	app := gin.New()
	app.Use(gin_logger.RecoveryWithZap(gin_logger.RecoveryLogger,
		gin_config.ServerConfigSettings.Server.IsDebug), gin_jaeger.JaegerMiddleWare(), cors.Default())
	return app
}

/*
	CreateServer 创建server
*/
func CreateServer(configPath, serverName string) (string, *gin.Engine) {
	initExtend(configPath, serverName)
	if !gin_config.ServerConfigSettings.Server.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	app := creatApp(configPath, serverName)
	endless.DefaultReadTimeOut = time.Duration(gin_config.ServerConfigSettings.Server.ReadTimeout) * time.Second
	endless.DefaultWriteTimeOut = time.Duration(gin_config.ServerConfigSettings.Server.WriteTimeout) * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf("%s:%d", gin_config.ServerConfigSettings.Server.ServerHost,
		gin_config.ServerConfigSettings.Server.ServerPort)
	return endPoint, app
}

/*
	DeferClose 关闭链接
*/
func DeferClose() {
	gin_mysql.CloseMysqlConnect()
	gin_jaeger.IoCloser()
	gin_redis.CloseRedisPool()
}
