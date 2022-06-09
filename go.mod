module github.com/fellowme/gin_common_library

go 1.15

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/alibabacloud-go/darabonba-openapi v0.1.14
	github.com/alibabacloud-go/dysmsapi-20170525/v2 v2.0.8
	github.com/alibabacloud-go/tea v1.1.17
	github.com/allegro/bigcache/v3 v3.0.2
	github.com/apache/pulsar-client-go v0.8.0
	github.com/fsnotify/fsnotify v1.5.1
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gomodule/redigo v1.8.5
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/olivere/elastic/v7 v7.0.29
	github.com/opentracing/opentracing-go v1.2.0
	github.com/panjf2000/ants/v2 v2.4.8
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.8.1
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/zap v1.17.0
	google.golang.org/grpc v1.38.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.3.4
	gorm.io/gorm v1.23.5
	gorm.io/plugin/soft_delete v1.1.0
)

replace github.com/spf13/viper v1.8.1 => github.com/spf13/viper v1.6.3
