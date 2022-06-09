package app

import (
	"github.com/fellowme/gin_common_library/big_cache"
	gin_config "github.com/fellowme/gin_common_library/config"
	gin_es "github.com/fellowme/gin_common_library/elastic"
	gin_jaeger "github.com/fellowme/gin_common_library/jaeger"
	gin_logger "github.com/fellowme/gin_common_library/logger"
	gin_pulsar "github.com/fellowme/gin_common_library/mq"
	gin_mysql "github.com/fellowme/gin_common_library/mysql"
	gin_redis "github.com/fellowme/gin_common_library/redis"
	gin_translator "github.com/fellowme/gin_common_library/translator"
	gin_util "github.com/fellowme/gin_common_library/util"
	"github.com/gin-gonic/gin"
	"net/http"
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
	gin_pulsar.InitPulsarClient()
	gin_es.InitElastic()
	big_cache.NewBigCache()
}

/*
	deferClose 关闭链接
*/
func deferClose() {
	gin_mysql.CloseMysqlConnect()
	gin_jaeger.IoCloser()
	gin_redis.CloseRedisPool()
	gin_pulsar.ClosePulsarClient()
}

/*
	initRouter  初始化路由
*/

func initRouter(app *gin.Engine, f func(group *gin.RouterGroup)) {
	// 测活
	api := app.Group("/api/v1")
	api.HEAD("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})
	if f != nil {
		f(api)
	}
}
