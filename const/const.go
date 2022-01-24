package _const

import (
	"time"
)

const (
	TokenNotEmptyTip    = "token 不能为空"
	TokenExpiredTip     = "token已过期"
	TokenInvalidTip     = "无法解析token"
	TokenMalformedTip   = "token不合法"
	TokenNotValidYetTip = "token未验证"
)

const (
	DefaultTxContextTimeOut = 3 * time.Second
	DefaultJwtExpiresAt     = 2 * time.Hour
)

const (
	MysqlNameDefault     = "default"
	DefaultMenuTableName = "gin_menu"
)

const (
	TimeFormat     = "2006-01-02 15:04:05"
	TimeFormatDate = "2006-01-02"
)
