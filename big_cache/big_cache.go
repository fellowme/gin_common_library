package big_cache

import (
	"github.com/allegro/bigcache/v3"
	gin_config "github.com/fellowme/gin_common_library/config"
	"go.uber.org/zap"
	"time"
)

var bigCache *bigcache.BigCache

func NewBigCache() {
	if gin_config.ServerConfigSettings.BigCacheConfig.Shards != 0 {
		config := bigcache.Config{Shards: gin_config.ServerConfigSettings.BigCacheConfig.Shards,
			LifeWindow:       gin_config.ServerConfigSettings.BigCacheConfig.LifeWindow * time.Minute,
			HardMaxCacheSize: gin_config.ServerConfigSettings.BigCacheConfig.HardMaxCacheSize,
			CleanWindow:      gin_config.ServerConfigSettings.BigCacheConfig.CleanWindow * time.Minute}
		var err error
		bigCache, err = bigcache.NewBigCache(config)
		if err != nil {
			zap.L().Error("NewBigCache error", zap.Any("error", err), zap.Any("config", gin_config.ServerConfigSettings.BigCacheConfig))
			return
		}
	}

}

func CloseBigCache() {
	if bigCache != nil {
		err := bigCache.Close()
		if err != nil {
			zap.L().Error("close bigCache error", zap.Any("error", err))
			return
		}
	}

}
