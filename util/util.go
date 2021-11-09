package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

const (
	SuccessCode = 0
	FailCode    = -1
)

var msgFlags = map[int]string{
	SuccessCode: "SUCCESS",
	FailCode:    "FAIL",
}

func getMsg(code int) string {
	msg, ok := msgFlags[code]
	if ok {
		return msg
	}
	return msgFlags[-1]
}

func ReturnResponse(httpCode, errCode int, message interface{}, data interface{}, c *gin.Context, version ...int) {
	c.JSON(httpCode, gin.H{
		"error_code":      errCode,
		"error_tip":       getMsg(errCode),
		"message":         message,
		"data":            data,
		"mapi_query_time": NowTimeToString(),
		"version":         version,
	})
}

//NowTimeToString *******现在日期格式化 yyyy-MM-dd HH:mm:ss*******//
func NowTimeToString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//RemoveSliceEmpty *******删除slice里面的空格和去除字符的前后的空格*******//
func RemoveSliceEmpty(arg []string) []string {
	var keys []string
	for _, key := range arg {
		if key != "" || strings.Trim(key, " ") != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

func Md5(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))

}
