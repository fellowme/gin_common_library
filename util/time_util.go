package util

import (
	gin_const "github.com/fellowme/gin_common_library/const"
	"time"
)

//NowTimeToString *******现在日期格式化 yyyy-MM-dd HH:mm:ss*******//
func NowTimeToString() string {
	return time.Now().Format(gin_const.TimeFormat)
}

func NowDateToString() string {
	return time.Now().Format(gin_const.TimeFormatDate)
}
