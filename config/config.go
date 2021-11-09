package config

import "time"

type ServerConfig struct {
	RedisConfigs  []RedisConf   `json:"redis_configs" form:"redis_configs"`
	MysqlConfigs  []MysqlConf   `json:"mysql_configs" form:"mysql_configs"`
	RocketMqConf  RocketMqConf  `json:"rocket_mq_conf" form:"rocket_mq_conf"`
	LoggerConfig  LoggerConfig  `json:"logger_config" form:"logger_config"`
	Server        Server        `json:"server" form:"server"`
	JaegerConfig  JaegerConfig  `json:"jaeger_config" form:"jaeger_config"`
	ElasticConfig ElasticConfig `json:"elastic_config" form:"elastic_config"`
}

type Server struct {
	ServerHost         string `json:"server_host" form:"server_host"`
	ServerPort         int    `json:"server_port" form:"server_port"`
	ServerName         string `json:"server_name" form:"server_name"`
	Path               string `json:"path" form:"path"`
	IsDebug            bool   `json:"is_debug" form:"is_debug"`
	PassportUrl        string `json:"passport_url" form:"passport_url"`
	RedisPrefix        string `json:"redis_prefix" form:"redis_prefix"`
	RedisCharacterMark string `json:"redis_character_mark" form:"redis_character_mark"`
}

type RocketMqConf struct {
	Endpoint    string            `json:"endpoint" form:"endpoint"`     // 设置HTTP接入域名（此处以公共云生产环境为例）
	AccessKey   string            `json:"access_key" form:"access_key"` // AccessKey 阿里云身份验证，在阿里云服务器管理控制台创建
	SecretKey   string            `json:"secret_key" form:"secret_key"` // SecretKey 阿里云身份验证，在阿里云服务器管理控制台创建
	InstanceIds map[string]string `json:"instance_ids" form:"instance_ids"`
	InstanceId  string            `json:"instance_id" form:"instance_id"` // Topic所属实例ID，默认实例为空
}

type RedisConf struct {
	Name             string        `json:"name" form:"name"`
	Host             string        `json:"host" form:"host"`
	Port             int           `json:"port" form:"port"`
	Password         string        `json:"password" form:"password"`
	Database         int           `json:"database" form:"database"`
	ConnectTimeout   time.Duration `json:"connect_timeout" form:"connect_timeout"`
	ReadTimeout      time.Duration `json:"read_timeout" form:"read_timeout"`
	ReadWriteTimeout time.Duration `json:"read_write_timeout" form:"read_write_timeout"`
	MaxIdle          int           `json:"max_idle" form:"max_idle"`
	MaxActive        int           `json:"max_active" form:"max_active"`
	IdleTimeout      time.Duration `json:"idle_timeout" form:"idle_timeout"`
	Wait             bool          `json:"wait" form:"wait"`
}

type MysqlConf struct {
	Name              string        `json:"name" form:"name"`
	Host              string        `json:"host" form:"host"`
	Port              string        `json:"port" form:"port"`
	User              string        `json:"user" form:"user"`
	Password          string        `json:"password" form:"password"`
	Database          string        `json:"database" form:"database"`
	MaxIdleConnects   int           `json:"max_idle_connects" form:"max_idle_connects"`
	MaxOpenConnects   int           `json:"max_open_connects" form:"max_open_connects"`
	ConnMaxLifetime   time.Duration `json:"conn_max_lifetime" form:"conn_max_lifetime"`
	LogModeBool       bool          `json:"log_mode_bool" form:"log_mode_bool"`
	SingularTableBool bool          `json:"singular_table_bool" form:"singular_table_bool"`
}

type LoggerConfig struct {
	LoggerPath       string `json:"logger_path" form:"logger_path"`
	LoggerMaxSize    int    `json:"logger_max_size" form:"logger_max_size"`
	LoggerMaxBackups int    `json:"logger_max_backups" form:"logger_max_backups"`
	LoggerMaxAge     int    `json:"logger_max_age" form:"logger_max_age"`
	LoggerIsCompress bool   `json:"logger_is_compress" form:"logger_is_compress"`
	LoggerLevelInt   int    `json:"logger_level_int" form:"logger_level_int"`
}

type JaegerConfig struct {
	Host  string  `json:"host" form:"host"`
	Port  int     `json:"port" form:"port"`
	Type  string  `json:"type" form:"type"`
	Param float64 `json:"param" form:"param"`
}

type ElasticConfig struct {
	Host     string `json:"host" form:"host"`
	Port     int    `json:"port" form:"port"`
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}
