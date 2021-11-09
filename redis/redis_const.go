package redis

const (
	ScriptLock = ` 
    local res=redis.call('GET', KEYS[1])
    if res then
        return 0
    else
        redis.call('SET',KEYS[1],ARGV[1]);
        redis.call('EXPIRE',KEYS[1],ARGV[2])
        return 1
    end 
    `

	ScriptExpire = ` 
    local res=redis.call('GET', KEYS[1])
    if not res then
        return -1
    end 
    if res==ARGV[1] then
        redis.call('EXPIRE', KEYS[1], ARGV[2])
        return 1
    else
        return 0
    end 
    `

	ScriptDelete = ` 
    local res=redis.call('GET', KEYS[1])
    if not res then 
        return -1
    end 
    if res==ARGV[1] then
        redis.call('DEL', KEYS[1])
    else
        return 0
    end 
    `
)

const (
	tcpConnect                = "tcp"
	redisDefault              = "default"
	retryCount                = 2
	redisPreFixDefault        = "default"
	redisSleepTime            = 10
	redisCharacterMarkDefault = "-"
)
