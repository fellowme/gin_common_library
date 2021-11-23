package util

import "time"

//NowTimeToString *******现在日期格式化 yyyy-MM-dd HH:mm:ss*******//
func NowTimeToString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
