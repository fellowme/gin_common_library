package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
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

func Md5(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))

}

func ReflectValueByType(value string, valueType string) interface{} {
	switch valueType {
	case "int":
		result, _ := strconv.Atoi(value)
		return result
	case "bool":
		result, _ := strconv.ParseBool(value)
		return result
	case "float64":
		result, _ := strconv.ParseFloat(value, 64)
		return result
	case "float32":
		result, _ := strconv.ParseFloat(value, 32)
		return result
	}

	return value
}

func GetPath() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}

	return path
}
