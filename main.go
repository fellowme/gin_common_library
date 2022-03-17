package main

import (
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	gin_logger "github.com/fellowme/gin_common_library/logger"
	gin_redis "github.com/fellowme/gin_common_library/redis"
)

func main() {
	gin_config.InitConfig("/Users/yanyongjun/go/src/gin_common_library", "go-library")
	gin_logger.InitServerLogger("/Users/yanyongjun/go/src/gin_common_library")
	gin_logger.InitRecoveryLogger("/Users/yanyongjun/go/src/gin_common_library")
	gin_redis.InitRedis()
	//code, err := gin_redis.SendScrip(gin_redis.ScriptIncrby, 1)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(code)
	//flag, err := gin_redis.SetRedisLockKey("11", "product_stock_lock", 800)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(flag)
	flag, err := gin_redis.DeleteRedisLockKey("11", "product_stock_lock")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(flag)
	//flag, err := gin_redis.ResetExpireRedisLockKey("11", "product_stock_lock", 800)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(flag)
}
