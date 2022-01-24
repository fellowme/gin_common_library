package router

import (
	gin_const "github.com/fellowme/gin_common_library/const"
	gin_mysql "github.com/fellowme/gin_common_library/mysql"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
	RegisterRouter 将router 注册到 menu
*/
func RegisterRouter(routers gin.RoutesInfo, mysqlDataBase ...string) {
	for _, router := range routers {
		routerStruct := struct {
			Method  string `json:"method"`
			Path    string `json:"path" `
			Handler string `json:"handler"`
		}{
			Method:  router.Method,
			Path:    router.Path,
			Handler: router.Handler,
		}
		err := gin_mysql.UseMysql(nil, mysqlDataBase...).Table(gin_const.DefaultMenuTableName).
			Where("is_delete =false and method=? and path=? and handler=?", router.Method,
				router.Path, router.Handler).FirstOrCreate(&routerStruct).Error
		if err != nil {
			zap.L().Error("menu registerRouter error", zap.Any("error", err.Error),
				zap.Any("router", router))
		}
	}
}
