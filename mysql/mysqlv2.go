package mysql

import (
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	gin_logger "github.com/fellowme/gin_common_library/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var mysqlV2Map map[string]*gorm.DB

func InitMysqlV2Map() {
	mysqlConfigs := gin_config.ServerConfigSettings.MysqlConfigs
	if len(mysqlConfigs) != 0 {
		mysqlV2Map = make(map[string]*gorm.DB, len(mysqlConfigs))
		for _, mysqlConfig := range mysqlConfigs {
			db := initMysqlV2(mysqlConfig)
			if db != nil {
				mysqlV2Map[mysqlConfig.Name] = db
			}
		}
	}
}

func initMysqlV2(mysqlConfig gin_config.MysqlConf) *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true",
		mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database)
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{Logger: gin_logger.NewSqlLogger(zap.L(), logger.Config{
		SlowThreshold:             200 * time.Nanosecond,
		Colorful:                  true,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  0,
	})})
	if err != nil {
		zap.L().Error("mysqlV2 open fail", zap.Any("error", err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("mysqlV2 sql.DB fail", zap.Any("error", err))
	}
	sqlDB.SetConnMaxLifetime(mysqlConfig.ConnMaxLifetime * time.Second)
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConnects)
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConnects)
	return db
}

func GetMysqlV2Map() map[string]*gorm.DB {
	if len(mysqlV2Map) == 0 {
		InitMysqlV2Map()
	}
	return mysqlV2Map
}

func UseMysqlV2(target map[string]*gorm.DB, name ...string) *gorm.DB {
	if target == nil {
		target = mysqlV2Map
	}
	mysqlName := mysqlNameDefault
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
