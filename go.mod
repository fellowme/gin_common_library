module github.com/fellowme/gin_commom_library

go 1.15

require (
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/jinzhu/gorm v1.9.16
	github.com/spf13/viper v1.9.0
	go.uber.org/zap v1.10.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/spf13/viper v1.9.0 => github.com/spf13/viper v1.6.3
