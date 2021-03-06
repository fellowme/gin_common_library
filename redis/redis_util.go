package redis

import (
	"context"
	gin_config "github.com/fellowme/gin_common_library/config"
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

// GetKey ****   获取name:redis名称 key:查询的key  返回 interface{} *****//
func GetKey(ctx context.Context, name, key string) (interface{}, error) {
	newKey := createKey(key)
	result, err := commandRedisWithRetry(ctx, name, "GET", newKey)
	return result, err
}

// MGetKey ****   获取name:redis名称 keys:查询的key  返回 interface{} *****//
func MGetKey(ctx context.Context, name string, keys []string) (interface{}, error) {
	keyList := make([]interface{}, 0)
	for _, key := range keys {
		keyList = append(keyList, createKey(key))
	}
	result, err := commandRedisWithRetry(ctx, name, "MGET", keyList...)
	return result, err
}

// MSetKey ****   设置name:redis名称 keys:设置的key value  返回 error *****//
func MSetKey(ctx context.Context, name string, keys []string) error {
	keyList := make([]interface{}, 0)
	for index, key := range keys {
		if index%2 == 0 {
			keyList = append(keyList, createKey(key))
		} else {
			keyList = append(keyList, key)
		}
	}
	_, err := commandRedisWithRetry(ctx, name, "MSET", keyList...)
	return err
}

// SetKeyValue  ****   设置name:redis名称 key:查询的key value:设定的值 expire:存在时间  返回 error *****//
func SetKeyValue(ctx context.Context, name, key string, value interface{}, expire ...int) (err error) {
	newKey := createKey(key)
	if len(expire) != 0 {
		_, err = commandRedisWithRetry(ctx, name, "SETEX", newKey, expire[0], value)
	} else {
		_, err = commandRedisWithRetry(ctx, name, "SET", newKey, value)
	}
	return
}

// ExistsKey  ****   是否存在 name:redis名称 key:查询的key  返回 bool error *****//
func ExistsKey(ctx context.Context, name, key string) (bool, error) {
	newKey := createKey(key)
	flag, err := redis.Bool(commandRedis(ctx, name, "EXISTS", newKey))
	return flag, err
}

// TtlKey  ****   返回 key 的剩余过期时间 name:redis名称 key:查询的key  返回 int64 error *****//
func TtlKey(ctx context.Context, name, key string) (int64, error) {
	newKey := createKey(key)
	ttl, err := redis.Int64(commandRedis(ctx, name, "TTL", newKey))
	return ttl, err
}

// ExpireKey  ****   设置key过期时间 name:redis名称 key:查询的key  返回 error *****//
func ExpireKey(ctx context.Context, name, key string, expire int64) error {
	newKey := createKey(key)
	_, err := redis.Int64(commandRedis(ctx, name, "EXPIRE", newKey, expire))
	return err
}

// IncrbyKey  ****   设置key增加指定的数值 name:redis名称 key:查询的key number:数值 返回 num error *****//
func IncrbyKey(ctx context.Context, name, key string, number int64) (num int64, err error) {
	newKey := createKey(key)
	num, err = redis.Int64(commandRedisWithRetry(ctx, name, "INCRBY", newKey, number))
	return
}

// DecrbyKey  ****   设置key减少指定的数值 name:redis名称 key:查询的key number:数值 返回 num error *****//
func DecrbyKey(ctx context.Context, name, key string, number int64) (num int64, err error) {
	newKey := createKey(key)
	num, err = redis.Int64(commandRedisWithRetry(ctx, name, "DECRBY", newKey, number))
	return
}

// DeleteKey  ****   删除key name:redis名称 key:删除key 返回 error *****//
func DeleteKey(ctx context.Context, name string, key string) error {
	newKey := createKey(key)
	_, err := commandRedisWithRetry(ctx, name, "DEL", newKey)
	return err
}

// SetBitmapKey  ****   设置bitmap  name:redis名称 key:删除key offset:偏移量 value:只能0，1  返回 error *****//
func SetBitmapKey(ctx context.Context, name, key string, offset, value int) error {
	newKey := createKey(key)
	_, err := commandRedis(ctx, name, "SETBIT", newKey, offset, value)
	if err != nil {
		zap.L().Error("SetBitmapKey fail", zap.String("newKey", newKey),
			zap.Int("offset", offset), zap.String("name", name), zap.String("key", key),
			zap.Int("value", value))
	}
	return err
}

// GetBitmapKey  ****   获取bitmap name:redis名称 key:删除key offset:偏移量  返回 int error *****//
func GetBitmapKey(ctx context.Context, name, key string, offset int) (int, error) {
	newKey := createKey(key)
	result, err := redis.Int(commandRedis(ctx, name, "GETBIT", newKey, offset))
	if err != nil {
		zap.L().Error("GetBitmapKey fail", zap.String("newKey", newKey),
			zap.Int("offset", offset), zap.String("name", name), zap.String("key", key))
	}
	return result, err
}

// SetRedisLockKey  ****   设置分布式锁 name:redis名称 key:值 lockName: 锁名字 time: 锁存在时间 返回 error *****//
func SetRedisLockKey(ctx context.Context, key, lockName string, time int, name ...string) (bool, error) {
	newKey := createKey(key)
	if ctx == nil {
		ctx = context.Background()
	}
	selectRedis, err := UseRedis(name...).GetContext(ctx)
	if err != nil {
		return false, err
	}
	luaExpire := redis.NewScript(1, ScriptLock)
	flag, err := redis.Bool(luaExpire.Do(selectRedis, lockName, newKey, time))
	if err != nil {
		zap.L().Error("SetRedisLockKey ScriptLock fail", zap.Any("error", err), zap.Any("name", name))
	}
	defer closeRedisConnect(selectRedis)
	return flag, err
}

// DeleteRedisLockKey  ****   删除分布式锁 name:redis名称 key:删除key lockName: 锁名字 返回 int *****//
func DeleteRedisLockKey(ctx context.Context, key, lockName string, name ...string) (bool, error) {
	newKey := createKey(key)
	if ctx == nil {
		ctx = context.Background()
	}
	selectRedis, err := UseRedis(name...).GetContext(ctx)
	if err != nil {
		return false, err
	}
	luaDel := redis.NewScript(1, ScriptDelete)
	flag, err := redis.Bool(luaDel.Do(selectRedis, lockName, newKey))
	if err != nil {
		zap.L().Error("DeleteRedisLockKey luaDel fail", zap.String("key", key), zap.Any("error", err), zap.Any("name", name))
	}
	defer closeRedisConnect(selectRedis)
	return flag, err
}

// ResetExpireRedisLockKey  ****   续费分布式锁 name:redis名称 key:删除key lockName: 锁名字 返回 int *****//
func ResetExpireRedisLockKey(ctx context.Context, key, lockName string, time int, name ...string) (bool, error) {
	newKey := createKey(key)
	if ctx == nil {
		ctx = context.Background()
	}
	selectRedis, err := UseRedis(name...).GetContext(ctx)
	if err != nil {
		return false, err
	}
	luaExpire := redis.NewScript(1, ScriptExpire)
	flag, err := redis.Bool(luaExpire.Do(selectRedis, lockName, newKey, time))
	if err != nil {
		zap.L().Error(" ResetExpireRedisLockKey luaExpire fail", zap.String("key", key), zap.Any("error", err), zap.Any("name", name))
	}
	defer closeRedisConnect(selectRedis)
	return flag, err
}

// ScriptDecrbyKey  ****   redis 脚本减少库存  *****//
func ScriptDecrbyKey(ctx context.Context, key string, number int, name ...string) (bool, error) {
	newKey := createKey(key)
	if ctx == nil {
		ctx = context.Background()
	}
	selectRedis, err := UseRedis(name...).GetContext(ctx)
	if err != nil {
		return false, err
	}
	luaScript := redis.NewScript(1, ScriptDecrby)
	flag, err := redis.Bool(luaScript.Do(selectRedis, newKey, number))
	if err != nil {
		zap.L().Error(" ScriptDecrbyKey NewScript do fail", zap.String("key", key), zap.Any("error", err), zap.Any("flag", flag))
	}
	defer closeRedisConnect(selectRedis)
	return flag, err
}

// ScriptIncrbyKey  ****   redis 脚本增加库存  *****//
func ScriptIncrbyKey(ctx context.Context, key string, number int, name ...string) (bool, error) {
	newKey := createKey(key)
	if ctx == nil {
		ctx = context.Background()
	}
	selectRedis, err := UseRedis(name...).GetContext(ctx)
	if err != nil {
		return false, err
	}
	luaScript := redis.NewScript(1, ScriptIncrby)
	flag, err := redis.Bool(luaScript.Do(selectRedis, newKey, number))
	if err != nil {
		zap.L().Error(" ScriptIncrbyKey NewScript do fail", zap.String("key", key), zap.Any("error", err), zap.Any("flag", flag))
	}
	defer closeRedisConnect(selectRedis)
	return flag, err
}

// SendScrip  ****   redis 上传脚本 不执行 *****//
func SendScrip(ctx context.Context, scriptString string, keyCount int, name ...string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	selectRedis, err := UseRedis(name...).GetContext(ctx)
	if err != nil {
		return "", err
	}
	luaExpire := redis.NewScript(keyCount, scriptString)
	err = luaExpire.Load(selectRedis)
	if err != nil {
		zap.L().Error(" SendScrip NewScript fail", zap.String("scriptString", scriptString), zap.Any("error", err))
	}
	hashCode := luaExpire.Hash()
	defer closeRedisConnect(selectRedis)
	return hashCode, err
}

// HSetMapKey 设置map值
func HSetMapKey(ctx context.Context, name, mapKey, key, value string) (err error) {
	newKey := createKey(mapKey)
	_, err = commandRedisWithRetry(ctx, name, "HSET", newKey, key, value)
	return
}

// HMSetMapKey 设置map值
func HMSetMapKey(ctx context.Context, name, mapKey string, args []interface{}) (err error) {
	newKey := createKey(mapKey)
	_, err = commandRedisWithRetry(ctx, name, "HMSET", append([]interface{}{
		newKey,
	}, args...)...)
	return
}

// HMGetMapKey 设置map值
func HMGetMapKey(ctx context.Context, name, mapKey string, args []interface{}) (data []string, err error) {
	newKey := createKey(mapKey)
	data, err = redis.Strings(commandRedisWithRetry(ctx, name, "HMGET", append([]interface{}{
		newKey,
	}, args...)...))
	return
}

// HGetMapKey 获取map值
func HGetMapKey(ctx context.Context, name, mapKey, key string) (value string, err error) {
	newKey := createKey(mapKey)
	value, err = redis.String(commandRedisWithRetry(ctx, name, "HGET", newKey, key))
	return
}

// HDelMapKey 获取map值
func HDelMapKey(ctx context.Context, name, mapKey, key string) (err error) {
	newKey := createKey(mapKey)
	_, err = commandRedisWithRetry(ctx, name, "HDEL", newKey, key)
	return
}
