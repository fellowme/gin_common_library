package middleware

import (
	"bytes"
	gin_logger "github.com/fellowme/gin_common_library/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/url"
	"time"
)

func RecordAccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		start := time.Now()
		proto := c.Request.Proto
		//修正特殊字符解析
		param, _ := url.QueryUnescape(c.Request.URL.RawQuery)
		c.Request.URL.RawQuery = param
		path := c.Request.URL.String()
		var bodyBytes []byte
		if c.Request.Body != nil {
			// 复制 request.body
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		cost := time.Since(start)
		list := []zap.Field{
			zap.String("ip", ip),
			zap.String("path", path),
			zap.String("proto", proto),
			zap.String("request_star_time", start.Format("2006-01-02 15:04:05")),
			zap.String("method", c.Request.Method),
			zap.String("request_param", string(bodyBytes)),
			zap.String("response_data", blw.body.String()),
			zap.Duration("cost", cost),
		}
		if len(c.Errors) > 0 {
			gin_logger.AccessLogger.Error(c.Errors.String(), list...)
		} else {
			gin_logger.AccessLogger.Info("gin success", list...)
		}

	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
