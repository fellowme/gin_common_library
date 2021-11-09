package gin_remote_service

import (
	"time"
)

const (
	// DefaultRetryTimes 如果请求失败，最多重试3次
	defaultRetryTimes = 3
	// DefaultRetryDelay 在重试前，延迟等待100毫秒
	defaultRetryDelay = time.Millisecond * 100
	// 默认ttl
	defaultTTl = 30 * time.Second
	// _StatusReadRespErr read resp body err, should re-call doHTTP again.
	_StatusReadRespErr = -204
	// _StatusDoReqErr do req err, should re-call doHTTP again.
	_StatusDoReqErr = -500
)
