package app

import (
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	gin_logger "github.com/fellowme/gin_common_library/logger"
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
func creatMqApp(configPath, serverName string) *gin.Engine {
	app := gin.New()
	app.Use(gin_logger.RecoveryWithZap(gin_logger.RecoveryLogger,
		gin_config.ServerConfigSettings.Server.IsDebug), cors.Default())
	return app
}

func CreateAppMqServer(configPath, serverName string, f func()) {
	// 执行配置
	initCommonExtend(configPath, serverName)
	// 是否执行 debug 模式
	if !gin_config.ServerConfigSettings.Server.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	//  初始化app
	app := creatMqApp(configPath, serverName)
	//  激活测活
	initRouter(app, nil)
	//  执行 方法
	f()
	//  endless 启动
	endless.DefaultReadTimeOut = time.Duration(gin_config.ServerConfigSettings.Server.ReadTimeout) * time.Second
	endless.DefaultWriteTimeOut = time.Duration(gin_config.ServerConfigSettings.Server.WriteTimeout) * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf("%s:%d", gin_config.ServerConfigSettings.Server.ServerHost,
		gin_config.ServerConfigSettings.Server.ServerMqPort)
	defer deferClose()
	server := endless.NewServer(endPoint, app)
	server.BeforeBegin = func(add string) {
		zap.L().Info(fmt.Sprintf("Actual pid is %d", syscall.Getpid()))
	}
	if err := server.ListenAndServe(); err != nil {
		panic(fmt.Sprint("init server fail err=", err))
	}
}
