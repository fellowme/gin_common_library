package big_cache

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

func FmtBigCacheKey(key string) string {
	return fmt.Sprintf(bigCacheKeyFmt, key)
}

func GetBigCacheValue(key string) (value []byte, err error) {
	newKey := FmtBigCacheKey(key)
	startTime := time.Now()
	zap.L().Info("GetBigCacheValue starting", zap.Any("key", key))
	value, err = bigCache.Get(newKey)
	if err != nil {
		zap.L().Error("GetBigCacheValue error", zap.Any("error", err), zap.Any("key", key))
		return
	}
	zap.L().Info("GetBigCacheValue end", zap.Any("key", key), zap.Any("value", string(value)), zap.Any("cost_time", time.Since(startTime).Seconds()))
	return
}

func SetBigCacheValue(key string, value []byte) (err error) {
	startTime := time.Now()
	newKey := FmtBigCacheKey(key)
	zap.L().Info("SetBigCacheValue starting", zap.Any("key", key), zap.Any("value", string(value)))
	err = bigCache.Set(newKey, value)
	if err != nil {
		zap.L().Error("SetBigCacheValue error", zap.Any("error", err), zap.Any("key", key), zap.Any("value", string(value)))
		return
	}
	zap.L().Info("SetBigCacheValue end", zap.Any("key", key), zap.Any("value", string(value)), zap.Any("cost_time", time.Since(startTime).Seconds()))
	return
}

func DeleteBigCacheValue(key string) (err error) {
	startTime := time.Now()
	newKey := FmtBigCacheKey(key)
	zap.L().Info("DeleteBigCacheValue starting", zap.Any("key", key))
	err = bigCache.Delete(newKey)
	if err != nil {
		zap.L().Error("DeleteBigCacheValue error", zap.Any("error", err), zap.Any("key", key))
		return
	}
	zap.L().Info("DeleteBigCacheValue end", zap.Any("key", key), zap.Any("cost_time", time.Since(startTime).Seconds()))
	return
}
