package app

import (
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	gin_jaeger "github.com/fellowme/gin_common_library/jaeger"
	gin_logger "github.com/fellowme/gin_common_library/logger"
	gin_mysql "github.com/fellowme/gin_common_library/mysql"
	gin_router "github.com/fellowme/gin_common_library/router"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"syscall"
	"time"
)

/*
	creatApp 初始化app
*/
func creatApp(configPath, serverName string) *gin.Engine {
	app := gin.New()
	app.Use(gin_logger.RecoveryWithZap(gin_logger.RecoveryLogger,
		gin_config.ServerConfigSettings.Server.IsDebug), gin_jaeger.TracerJaegerMiddleWare(), cors.Default())
	return app
}

func CreateAppServer(configPath, serverName string, f func(group *gin.RouterGroup), models []interface{}) {
	// 执行配置
	initCommonExtend(configPath, serverName)
	// 是否执行 debug 模式
	if !gin_config.ServerConfigSettings.Server.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	//  初始化app
	app := creatApp(configPath, serverName)
	//  初始化 model生成数据库表
	initTable(models)
	//  激活测活  加入新的路由
	initRouter(app, f)
	//  将新路由 添加到数据库
	gin_router.RegisterRouter(app.Routes(), serverName)
	//  endless 启动
	endless.DefaultReadTimeOut = time.Duration(gin_config.ServerConfigSettings.Server.ReadTimeout) * time.Second
	endless.DefaultWriteTimeOut = time.Duration(gin_config.ServerConfigSettings.Server.WriteTimeout) * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf("%s:%d", gin_config.ServerConfigSettings.Server.ServerHost,
		gin_config.ServerConfigSettings.Server.ServerPort)
	defer deferClose()
	server := endless.NewServer(endPoint, app)
	server.BeforeBegin = func(add string) {
		zap.L().Info(fmt.Sprintf("Actual pid is %d", syscall.Getpid()))
	}
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprint("init server fail err=", err))
	}
}

/*
	initTable 初始化mysql 表信息
*/
func initTable(models []interface{}) {
	if models != nil {
		err := gin_mysql.UseMysql(nil).AutoMigrate(models...)
		if err != nil {
			zap.L().Error("UseMysql error", zap.Any("error", err))
		}
	}
}
