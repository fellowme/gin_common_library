package mysql

import (
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"time"
)

const mysqlNameDefault = "default"

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

func GetMysqlMap() map[string]*gorm.DB {
	if len(mysqlMap) == 0 {
		InitMysqlMap()
	}
	return mysqlMap
}

func UseMysql(name ...string) *gorm.DB {
	mysqlName := mysqlNameDefault
	if len(name) != 0 {
		mysqlName = name[0]
	}
	mysql, ok := mysqlMap[mysqlName]
	if !ok {
		zap.L().Error("UseMysql fail not find right connect", zap.Any("name", name))
		return nil
	}
	return mysql
}

func initMysql(mysqlConfig gin_config.MysqlConf) *gorm.DB {
	var err error
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true",
		mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database)
	db, err := gorm.Open("mysql", url)
	if err != nil {
		zap.L().Error("mysql init error", zap.Any("error", err), zap.Any("mysqlConfig", mysqlConfig))
		return nil
	}
	db.LogMode(mysqlConfig.LogModeBool)
	db.SingularTable(mysqlConfig.SingularTableBool)
	db.DB().SetConnMaxLifetime(mysqlConfig.ConnMaxLifetime * time.Second)
	db.DB().SetMaxIdleConns(mysqlConfig.MaxIdleConnects)
	db.DB().SetMaxOpenConns(mysqlConfig.MaxOpenConnects)
	//db.Callback().Create().Before("gorm:create").Register("update_created_time", updateTimeStampForCreateCallback)
	//db.Callback().Create().Before("gorm:update").Register("update_time", updateTimeStampForUpdateCallback)
	return db
}

func CloseMysqlPool() {
	if len(mysqlMap) != 0 {
		for key, mysqlPool := range mysqlMap {
			if err := mysqlPool.Close(); err != nil {
				zap.L().Error("mysql close error", zap.String("mysqlName", key))
			}
		}
	}
}

//// updateTimeStampForCreateCallback 注册新建钩子在持久化之前  *******创建之前******//
//func updateTimeStampForCreateCallback(scope *gorm.Scope) {
//	if !scope.HasError() {
//		if createTimeField, ok := scope.FieldByName("create_time"); ok {
//			if createTimeField.IsBlank {
//				if err := createTimeField.Set(time.Now()); err != nil {
//					zap.L().Error("updateTimeStampForCreateCallback  createTimeField error", zap.Any("error", err))
//				}
//			}
//		}
//		if modifyTimeField, ok := scope.FieldByName("update_time"); ok {
//			if modifyTimeField.IsBlank {
//				if err := modifyTimeField.Set(time.Now()); err != nil {
//					zap.L().Error("updateTimeStampForCreateCallback  modifyTimeField error", zap.Any("error", err))
//				}
//			}
//		}
//	}
//}
//
//// updateTimeStampForUpdateCallback 注册新建钩子在持久化之前  *******更新之前******//
//func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
//	if !scope.HasError() {
//		if modifyTimeField, ok := scope.FieldByName("update_time"); ok {
//			if modifyTimeField.IsBlank {
//				if err := modifyTimeField.Set(time.Now()); err != nil {
//					zap.L().Error("updateTimeStampForCreateCallback  modifyTimeField error", zap.Any("error", err))
//				}
//			}
//		}
//	}
//}
