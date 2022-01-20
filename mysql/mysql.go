package mysql

import (
	"context"
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	gin_const "github.com/fellowme/gin_common_library/const"
	gin_logger "github.com/fellowme/gin_common_library/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var mysqlMap map[string]*gorm.DB

func InitMysqlMap() {
	mysqlConfigs := gin_config.ServerConfigSettings.MysqlConfigs
	if len(mysqlConfigs) != 0 {
		mysqlMap = make(map[string]*gorm.DB, len(mysqlConfigs))
		for _, mysqlConfig := range mysqlConfigs {
			db := initMysql(mysqlConfig)
			if db != nil {
				mysqlMap[mysqlConfig.Name] = db
			}
		}
	}
}

func initMysql(mysqlConfig gin_config.MysqlConf) *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true",
		mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{Logger: gin_logger.NewSqlLogger(zap.L(), logger.Config{
		SlowThreshold:             mysqlConfig.SlowThreshold * time.Second,
		Colorful:                  mysqlConfig.Colorful,
		IgnoreRecordNotFoundError: mysqlConfig.IgnoreRecordNotFoundError,
		LogLevel:                  logger.LogLevel(mysqlConfig.LogLevel),
	})})
	if err != nil {
		zap.L().Error("mysql open fail", zap.Any("error", err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("mysql sql.DB fail", zap.Any("error", err))
	}
	sqlDB.SetConnMaxLifetime(mysqlConfig.ConnMaxLifetime * time.Second)
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConnects)
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConnects)
	return db
}

func GetMysqlMap() map[string]*gorm.DB {
	if len(mysqlMap) == 0 {
		InitMysqlMap()
	}
	return mysqlMap
}

func UseMysql(target map[string]*gorm.DB, name ...string) *gorm.DB {
	if target == nil {
		target = mysqlMap
	}
	mysqlName := gin_const.MysqlNameDefault
	if len(name) != 0 {
		mysqlName = name[0]
	}
	mysqlDB, ok := target[mysqlName]
	if !ok {
		zap.L().Error("UseMysql fail not find right connect", zap.Any("name", name))
		return nil
	}
	return mysqlDB
}

func GetTxWithContext(target map[string]*gorm.DB, ctx context.Context, tableName string, name ...string) (*gorm.DB, context.CancelFunc) {
	mysqlDB := UseMysql(target, name...)
	if mysqlDB == nil {
		return nil, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	contextTimeout, cancel := context.WithTimeout(ctx, gin_const.DefaultTxContextTimeOut)
	tx := mysqlDB.WithContext(contextTimeout).Table(tableName)
	return tx, cancel
}
