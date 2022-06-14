package router

import (
	gin_const "github.com/fellowme/gin_common_library/const"
	gin_model "github.com/fellowme/gin_common_library/model"
	gin_mysql "github.com/fellowme/gin_common_library/mysql"
	gin_util "github.com/fellowme/gin_common_library/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

/*
	RegisterRouter 将router 注册到 menu
*/
func RegisterRouter(routers gin.RoutesInfo, serverName string, mysqlDataBase ...string) {
	if gin_mysql.UseMysql(nil, mysqlDataBase...) == nil {
		return
	}
	err := gin_mysql.UseMysql(nil, mysqlDataBase...).Table(gin_const.DefaultMenuTableName).
		Where("server_name = ?", serverName).Delete(&MenuStructParam{}).Error
	if err != nil {
		zap.L().Error("menu registerRouter Delete error", zap.Any("error", err.Error),
			zap.Any("serverName", serverName))
	}
	for _, router := range routers {
		// head 请求不加入
		if strings.ToLower(router.Method) != gin_util.RequestHeadMethod {
			routerStruct := MenuStructParam{
				Method:     router.Method,
				Path:       router.Path,
				Handler:    router.Handler,
				ServerName: serverName,
			}
			var total int64
			db := gin_mysql.UseMysql(nil, mysqlDataBase...).Table(gin_const.DefaultMenuTableName)
			db = db.Where("method = ? and path = ? and server_name = ?", router.Method,
				router.Path, serverName).Count(&total)
			if total > 0 {
				db.Updates(map[string]interface{}{
					"is_delete": 0,
				})
			} else {
				db.Create(&routerStruct)
			}
			err = db.Error
			if err != nil {
				zap.L().Error("menu registerRouter error", zap.Any("error", err.Error),
					zap.Any("router", router))
			}
		}

	}
}

type MenuStructParam struct {
	gin_model.BaseMysqlStruct
	Method     string `json:"method"`
	Path       string `json:"path" `
	Handler    string `json:"handler"`
	ServerName string `json:"server_name"`
}
