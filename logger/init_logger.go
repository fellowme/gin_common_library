package logger

import (
	gin_config "github.com/fellowme/gin_commom_library/config"
	"go.uber.org/zap"
)

var AccessLogger *zap.Logger
var RecoveryLogger *zap.Logger

func InitAccessLogger(basePath string) {
	config := gin_config.ServerConfigSettings.LoggerConfig
	serverName := basePath + config.LoggerPath + gin_config.ServerConfigSettings.Server.ServerName + "_access"
	AccessLogger = initLogger(config, serverName)
}

func InitServerLogger(basePath string) {
	config := gin_config.ServerConfigSettings.LoggerConfig
	serverName := basePath + config.LoggerPath + gin_config.ServerConfigSettings.Server.ServerName
	GinLogger := initLogger(config, serverName)
	zap.ReplaceGlobals(GinLogger)
}

func InitRecoveryLogger(basePath string) {
	config := gin_config.ServerConfigSettings.LoggerConfig
	serverName := basePath + config.LoggerPath + gin_config.ServerConfigSettings.Server.ServerName + "_recovery"
	RecoveryLogger = initLogger(config, serverName)
}
