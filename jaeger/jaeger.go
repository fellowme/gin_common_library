package jaeger

import (
	"bytes"
	"fmt"
	gin_config "github.com/fellowme/gin_commom_library/config"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
)

var tracer opentracing.Tracer
var closer io.Closer

func InitJaegerTracer() {
	var err error
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  gin_config.ServerConfigSettings.JaegerConfig.Type,  //百分比采样率
			Param: gin_config.ServerConfigSettings.JaegerConfig.Param, //按照百分比采样
		},

		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", gin_config.ServerConfigSettings.JaegerConfig.Host, gin_config.ServerConfigSettings.JaegerConfig.Port),
		},

		ServiceName: gin_config.ServerConfigSettings.Server.ServerName,
	}

	tracer, closer, err = cfg.NewTracer()
	if err != nil {
		zap.L().Error("InitJaegerTracer fail", zap.Any("error ", err))
	}
	opentracing.SetGlobalTracer(tracer)
}

func IoCloser() {
	if err := closer.Close(); err != nil {
		zap.L().Error("JaegerTracer IoCloser fail", zap.Any("error", err))
	}
}

func JaegerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var parentSpan opentracing.Span
		var bodyBytes []byte
		if c.Request.Body != nil {
			// 复制 request.body
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			parentSpan = tracer.StartSpan(c.Request.URL.Path,
				opentracing.Tags{
					"method":              c.Request.Method,
					"url":                 c.Request.URL,
					"data":                string(bodyBytes),
					"header":              c.Request.Header,
					string(ext.Component): c.Request.Proto,
					"serverName":          gin_config.ServerConfigSettings.Server.ServerName,
				},
				ext.SpanKindRPCServer)
			defer parentSpan.Finish()
		} else {
			parentSpan = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tags{
					"method":              c.Request.Method,
					"url":                 c.Request.URL,
					"data":                string(bodyBytes),
					"header":              c.Request.Header,
					string(ext.Component): c.Request.Proto,
					"serverName":          gin_config.ServerConfigSettings.Server.ServerName,
				},
				ext.SpanKindRPCServer,
			)
			defer parentSpan.Finish()
		}
		c.Set("Tracer", tracer)
		c.Set("ParentSpanContext", parentSpan.Context())
		c.Next()
	}
}
