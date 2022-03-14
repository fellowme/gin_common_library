package router

import (
	gin_const "github.com/fellowme/gin_common_library/const"
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
	for _, router := range routers {
		// head 请求不加入
		if strings.ToLower(router.Method) != gin_util.RequestHeadMethod {
			routerStruct := struct {
				Method     string `json:"method"`
				Path       string `json:"path" `
				Handler    string `json:"handler"`
				ServerName string `json:"server_name"`
			}{
				Method:     router.Method,
				Path:       router.Path,
				Handler:    router.Handler,
				ServerName: serverName,
			}
			err := gin_mysql.UseMysql(nil, mysqlDataBase...).Table(gin_const.DefaultMenuTableName).
				Where("is_delete = false and method = ? and path = ? and server_name = ?", router.Method,
					router.Path, router.Handler).FirstOrCreate(&routerStruct).Error
			if err != nil {
				zap.L().Error("menu registerRouter error", zap.Any("error", err.Error),
					zap.Any("router", router))
			}
		}

	}
}
