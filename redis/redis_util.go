package redis

import (
	"errors"
	gin_config "github.com/fellowme/gin_common_library/config"
	gin_util "github.com/fellowme/gin_common_library/util"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"strings"
)

func createKey(key string) string {
	if key == "" {
		return ""
	}
	redisPrefix := redisPreFixDefault
	redisCharacterMark := redisCharacterMarkDefault
	if gin_config.ServerConfigSettings.Server.RedisPrefix != "" {
		redisPrefix = gin_config.ServerConfigSettings.Server.RedisPrefix
	}
	if gin_config.ServerConfigSettings.Server.RedisCharacterMark != "" {
		redisCharacterMark = gin_config.ServerConfigSettings.Server.RedisCharacterMark
	}
	return strings.Join([]string{redisPrefix, key}, redisCharacterMark)
}

// GetKeyByte ****   获取name:redis名称 key:查询的key  返回 interface{} *****//
func GetKeyByte(name, key string) (interface{}, error) {
	newKey := createKey(key)
	result, err := commandRedisWithRetry(name, "GET", newKey)
	return result, err
}

// SetKeyValue  ****   设置name:redis名称 key:查询的key value:设定的值 expire:存在时间  返回 error *****//
func SetKeyValue(name, key string, value interface{}, expire ...int) (err error) {
	newKey := createKey(key)
	if len(expire) != 0 {
		_, err = commandRedisWithRetry(name, "SET", newKey, value, "EX", expire[0])
	} else {
		_, err = commandRedisWithRetry(name, "SET", newKey, value)
	}
	return
}

// ExistsKey  ****   是否存在 name:redis名称 key:查询的key  返回 bool error *****//
func ExistsKey(name, key string) (bool, error) {
	newKey := createKey(key)
	flag, err := redis.Bool(commandRedis(name, "EXISTS", newKey))
	return flag, err
}

// TtlKey  ****   返回 key 的剩余过期时间 name:redis名称 key:查询的key  返回 int64 error *****//
func TtlKey(name, key string) (int64, error) {
	newKey := createKey(key)
	ttl, err := redis.Int64(commandRedis(name, "TTL", newKey))
	return ttl, err
}

// ExpireKey  ****   设置key过期时间 name:redis名称 key:查询的key  返回 error *****//
func ExpireKey(name, key string, expire int64) error {
	newKey := createKey(key)
	_, err := redis.Int64(commandRedis(name, "EXPIRE", newKey, expire))
	return err
}

// IncrbyKey  ****   设置key增加指定的数值 name:redis名称 key:查询的key number:数值 返回 num error *****//
func IncrbyKey(name, key string, number int64) (num int64, err error) {
	newKey := createKey(key)
	num, err = redis.Int64(commandRedisWithRetry(name, "INCRBY", newKey, number))
	return
}

// DecrbyKey  ****   设置key减少指定的数值 name:redis名称 key:查询的key number:数值 返回 num error *****//
func DecrbyKey(name, key string, number int64) (num int64, err error) {
	newKey := createKey(key)
	num, err = redis.Int64(commandRedisWithRetry(name, "DECRBY", newKey, number))
	return
}

// DeleteKeys  ****   删除多个key name:redis名称 key:删除key 返回 error *****//
func DeleteKeys(name string, key ...string) error {
	if len(key) == 0 {
		return errors.New("delete key empty")
	}
	keys := gin_util.RemoveSliceEmpty(key)
	if len(keys) == 0 {
		return errors.New("delete RemoveSliceEmpty key empty")
	}
	newKeys := make([]string, 0)
	for _, key := range keys {
		newKeys = append(newKeys, createKey(key))
	}
	_, err := commandRedisWithRetry(name, "DEL", newKeys)
	return err
}

// SetBitmapKey  ****   设置bitmap  name:redis名称 key:删除key offset:偏移量 value:只能0，1  返回 error *****//
func SetBitmapKey(name, key string, offset, value int) error {
	newKey := createKey(key)
	_, err := commandRedis(name, "SetBit", newKey, offset, value)
	if err != nil {
		zap.L().Error("SetBitmapKey fail", zap.String("newKey", newKey),
			zap.Int("offset", offset), zap.String("name", name), zap.String("key", key),
			zap.Int("value", value))
	}
	return err
}

// GetBitmapKey  ****   获取bitmap name:redis名称 key:删除key offset:偏移量  返回 int error *****//
func GetBitmapKey(name, key string, offset int) (int, error) {
	newKey := createKey(key)
	result, err := redis.Int(commandRedis(name, "GetBit", newKey, offset))
	if err != nil {
		zap.L().Error("GetBitmapKey fail", zap.String("newKey", newKey),
			zap.Int("offset", offset), zap.String("name", name), zap.String("key", key))
	}
	return result, err
}

// SetRedisLockKey  ****   设置分布式锁 name:redis名称 key:删除key lockName: 锁名字 time: 锁存在时间 返回 error *****//
func SetRedisLockKey(name, key, lockName string, time int) error {
	newKey := createKey(key)
	selectRedis := UseRedis(name).Get()
	luaExpire := redis.NewScript(1, ScriptLock)
	_, err := redis.Int(luaExpire.Do(selectRedis, lockName, newKey, time))
	if err != nil {
		zap.L().Error("luaExpire fail", zap.String("key", key), zap.Any("error", err), zap.String("name", name))
	}
	defer closeRedisConnect(selectRedis)
	return err
}

// DeleteRedisLockKey  ****   删除分布式锁 name:redis名称 key:删除key lockName: 锁名字 返回 int *****//
func DeleteRedisLockKey(name, key, lockName string) error {
	newKey := createKey(key)
	selectRedis := UseRedis(name).Get()
	luaDel := redis.NewScript(1, ScriptDelete)
	_, err := redis.Int(luaDel.Do(selectRedis, lockName, newKey))
	if err != nil {
		zap.L().Error("luaDel fail", zap.String("key", key), zap.Any("error", err), zap.String("name", name))
	}
	defer closeRedisConnect(selectRedis)
	return err
}

// ResetExpireRedisLockKey  ****   续费分布式锁 name:redis名称 key:删除key lockName: 锁名字 返回 int *****//
func ResetExpireRedisLockKey(name, key, lockName string, time int) error {
	newKey := createKey(key)
	selectRedis := UseRedis(name).Get()
	luaExpire := redis.NewScript(1, ScriptExpire)
	_, err := redis.Int(luaExpire.Do(selectRedis, lockName, newKey, time))
	if err != nil {
		zap.L().Error(" ResetExpireRedisLockKey luaExpire fail", zap.String("key", key), zap.Any("error", err), zap.String("name", name))
	}
	defer closeRedisConnect(selectRedis)
	return err
}
