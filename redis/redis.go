package redis

import (
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"time"
)

var redisMap map[string]*redis.Pool

func InitRedis() {
	if len(gin_config.ServerConfigSettings.RedisConfigs) > 0 {
		redisMap = make(map[string]*redis.Pool, len(gin_config.ServerConfigSettings.RedisConfigs))
		for _, redisConfig := range gin_config.ServerConfigSettings.RedisConfigs {
			redisPool := &redis.Pool{
				Dial:        GetRedisConnect(redisConfig),
				MaxIdle:     redisConfig.MaxIdle,
				MaxActive:   redisConfig.MaxActive,
				IdleTimeout: redisConfig.IdleTimeout * time.Second,
				Wait:        redisConfig.Wait,
			}
			redisMap[redisConfig.Name] = redisPool
		}
	}
}

func GetRedisConnect(redisConfig gin_config.RedisConf) func() (redis.Conn, error) {
	client, err := redis.Dial(tcpConnect, fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		redis.DialConnectTimeout(redisConfig.ConnectTimeout*time.Second), redis.DialReadTimeout(redisConfig.ReadTimeout*time.Second),
		redis.DialWriteTimeout(redisConfig.ReadWriteTimeout*time.Second))
	if err != nil {
		zap.L().Error("ERROR: fail init redis pool", zap.Any("error", err))
	}
	if redisConfig.Password != "" {
		if _, err = client.Do("AUTH", redisConfig.Password); err != nil {
			closeRedisConnect(client)
			zap.L().Error("ERROR: fail Password redis pool", zap.Any("error", err))
			return func() (redis.Conn, error) {
				return nil, err
			}
		}
	}
	if redisConfig.Database > 0 {
		if _, err = client.Do("SELECT", redisConfig.Database); err != nil {
			zap.L().Error("ERROR: fail SELECT redis pool", zap.Any("error", err))
			closeRedisConnect(client)
			return func() (redis.Conn, error) {
				return nil, err
			}
		}
	}
	return func() (redis.Conn, error) {
		return client, err
	}
}

// UseRedis 获取使用的redis
func UseRedis(name ...string) *redis.Pool {
	if len(redisMap) == 0 {
		InitRedis()
	}
	redisName := redisDefault
	if len(name) != 0 {
		redisName = name[0]
	}
	selectRedis, ok := redisMap[redisName]
	if !ok {
		panic("无法获取 redis_" + redisName)
	}
	return selectRedis
}

// redis 重试2次  间隔 10ms
func commandRedisWithRetry(name, commandName string, args ...interface{}) (reply interface{}, err error) {
	for i := 0; i < retryCount; i++ {
		reply, err = commandRedis(name, commandName, args...)
		if err == nil {
			return
		}
		time.Sleep(redisSleepTime * time.Millisecond)
	}
	return
}

// 执行redis command 命令
func commandRedis(name, commandName string, args ...interface{}) (reply interface{}, err error) {
	selectRedis := UseRedis(name).Get()
	nowTime := time.Now()
	reply, err = selectRedis.Do(commandName, args...)
	costTime := time.Since(nowTime)
	if err != nil {
		zap.L().Error("redis do fail", zap.Any("error", err), zap.Duration("costTime", costTime),
			zap.String("commandName", commandName), zap.Any("args", args))
	} else {
		zap.L().Info("redis do success", zap.Duration("costTime", costTime),
			zap.String("commandName", commandName), zap.Any("args", args))

	}
	defer closeRedisConnect(selectRedis)
	return
}

// 归还redis pool connect
func closeRedisConnect(redisConnect redis.Conn) {
	if err := redisConnect.Close(); err != nil {
		zap.L().Error("redis connect close fail", zap.Any("error", err))
	}
}

// CloseRedisPool 关闭 所有redis pool
func CloseRedisPool() {
	if len(redisMap) != 0 {
		for index, pool := range redisMap {
			if err := pool.Close(); err != nil {
				zap.L().Error("redis pool close fail", zap.Any("error", err), zap.String("redis name", index))
				continue
			}
		}
	}
}
