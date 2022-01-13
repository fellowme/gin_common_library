module github.com/fellowme/gin_common_library

go 1.15

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/alibabacloud-go/darabonba-openapi v0.1.14
	github.com/alibabacloud-go/dysmsapi-20170525/v2 v2.0.8
	github.com/alibabacloud-go/tea v1.1.17
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/gomodule/redigo v1.8.5
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/jinzhu/gorm v1.9.16
	github.com/olivere/elastic/v7 v7.0.29
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.9.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/zap v1.10.0
	google.golang.org/grpc v1.33.2
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.2.1
	gorm.io/gorm v1.22.4
	gorm.io/plugin/soft_delete v1.0.5
)

replace github.com/spf13/viper v1.9.0 => github.com/spf13/viper v1.6.3
