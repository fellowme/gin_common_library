package config

import "time"

type ServerConfig struct {
	RedisConfigs              []RedisConf               `json:"redis_configs" form:"redis_configs" mapstructure:"redis_configs"`
	MysqlConfigs              []MysqlConf               `json:"mysql_configs" form:"mysql_configs" mapstructure:"mysql_configs"`
	PulsarMqConf              PulsarMqConf              `json:"pulsar_mq_conf" mapstructure:"pulsar_mq_conf"`
	LoggerConfig              LoggerConfig              `json:"logger_config" form:"logger_config" mapstructure:"logger_config"`
	Server                    Server                    `json:"server" form:"server" mapstructure:"server"`
	JaegerConfig              JaegerConfig              `json:"jaeger_config" form:"jaeger_config" mapstructure:"jaeger_config"`
	ElasticConfig             ElasticConfig             `json:"elastic_config" form:"elastic_config" mapstructure:"elastic_config"`
	AliYunSendPhoneCodeConfig AliYunSendPhoneCodeConfig `json:"ali_yun_send_phone_code_config" mapstructure:"ali_yun_send_phone_code_config"`
}

type Server struct {
	SignKey            string `json:"sign_key" mapstructure:"sign_key"`
	ReadTimeout        int    `json:"read_timeout" form:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout       int    `json:"write_timeout" form:"write_timeout" mapstructure:"write_timeout"`
	ServerHost         string `json:"server_host" form:"server_host" mapstructure:"server_host"`
	ServerPort         int    `json:"server_port" form:"server_port" mapstructure:"server_port"`
	ServerRpcPort      int    `json:"server_rpc_port" form:"server_rpc_port" mapstructure:"server_rpc_port"`
	ServerMqPort       int    `json:"server_mq_port" form:"server_mq_port" mapstructure:"server_mq_port"`
	ServerName         string `json:"server_name" form:"server_name" mapstructure:"server_name"`
	Path               string `json:"path" form:"path" mapstructure:"path"`
	IsDebug            bool   `json:"is_debug" form:"is_debug" mapstructure:"is_debug"`
	RedisPrefix        string `json:"redis_prefix" form:"redis_prefix" mapstructure:"redis_prefix"`
	RedisCharacterMark string `json:"redis_character_mark" form:"redis_character_mark" mapstructure:"redis_character_mark"`
}

type RedisConf struct {
	Name             string        `json:"name" form:"name" mapstructure:"name"`
	Host             string        `json:"host" form:"host" mapstructure:"host"`
	Port             int           `json:"port" form:"port" mapstructure:"port"`
	Password         string        `json:"password" form:"password" mapstructure:"password"`
	Database         int           `json:"database" form:"database" mapstructure:"database"`
	ConnectTimeout   time.Duration `json:"connect_timeout" form:"connect_timeout" mapstructure:"connect_timeout"`
	ReadTimeout      time.Duration `json:"read_timeout" form:"read_timeout" mapstructure:"read_timeout"`
	ReadWriteTimeout time.Duration `json:"read_write_timeout" form:"read_write_timeout" mapstructure:"read_write_timeout"`
	MaxIdle          int           `json:"max_idle" form:"max_idle" mapstructure:"max_idle"`
	MaxActive        int           `json:"max_active" form:"max_active" mapstructure:"max_active"`
	IdleTimeout      time.Duration `json:"idle_timeout" form:"idle_timeout" mapstructure:"idle_timeout"`
	Wait             bool          `json:"wait" form:"wait" mapstructure:"wait"`
}

type MysqlConf struct {
	Name                      string        `json:"name" form:"name" mapstructure:"name"`
	Host                      string        `json:"host" form:"host" mapstructure:"host"`
	Port                      string        `json:"port" form:"port" mapstructure:"port"`
	User                      string        `json:"user" form:"user" mapstructure:"user"`
	Password                  string        `json:"password" form:"password" mapstructure:"password"`
	Database                  string        `json:"database" form:"database" mapstructure:"database"`
	MaxIdleConnects           int           `json:"max_idle_connects" form:"max_idle_connects" mapstructure:"max_idle_connects"`
	MaxOpenConnects           int           `json:"max_open_connects" form:"max_open_connects" mapstructure:"max_open_connects"`
	ConnMaxLifetime           time.Duration `json:"conn_max_lifetime" form:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	LogModeBool               bool          `json:"log_mode_bool" form:"log_mode_bool" mapstructure:"log_mode_bool"`
	SingularTableBool         bool          `json:"singular_table_bool" form:"singular_table_bool" mapstructure:"singular_table_bool"`
	Colorful                  bool          `json:"colorful" mapstructure:"colorful"`
	IgnoreRecordNotFoundError bool          `json:"ignore_record_not_found_error" mapstructure:"ignore_record_not_found_error"`
	SlowThreshold             time.Duration `json:"slow_threshold" mapstructure:"slow_threshold"`
	LogLevel                  int           `json:"log_level" mapstructure:"log_level"`
}

type LoggerConfig struct {
	LoggerPath       string `json:"logger_path" form:"logger_path" mapstructure:"logger_path"`
	LoggerMaxSize    int    `json:"logger_max_size" form:"logger_max_size" mapstructure:"logger_max_size"`
	LoggerMaxBackups int    `json:"logger_max_backups" form:"logger_max_backups" mapstructure:"logger_max_backups"`
	LoggerMaxAge     int    `json:"logger_max_age" form:"logger_max_age" mapstructure:"logger_max_age"`
	LoggerIsCompress bool   `json:"logger_is_compress" form:"logger_is_compress" mapstructure:"logger_is_compress"`
	LoggerLevelInt   int    `json:"logger_level_int" form:"logger_level_int" mapstructure:"logger_level_int"`
}

type JaegerConfig struct {
	Host  string  `json:"host" form:"host" mapstructure:"host"`
	Port  int     `json:"port" form:"port" mapstructure:"port"`
	Type  string  `json:"type" form:"type" mapstructure:"type"`
	Param float64 `json:"param" form:"param" mapstructure:"param"`
}

type ElasticConfig struct {
	Host     string `json:"host" form:"host" mapstructure:"host"`
	Port     int    `json:"port" form:"port" mapstructure:"port"`
	UserName string `json:"user_name" form:"user_name" mapstructure:"user_name"`
	Password string `json:"password" form:"password" mapstructure:"password"`
}

type AliYunSendPhoneCodeConfig struct {
	AccessKeyId     string `json:"access_key_id,omitempty" mapstructure:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret,omitempty" mapstructure:"access_key_secret"`
	Endpoint        string `json:"endpoint" form:"endpoint" mapstructure:"endpoint"`
}

type PulsarMqConf struct {
	PulsarUrl           string            `json:"pulsar_url" mapstructure:"pulsar_url"`                 // 设置HTTP接入域名（此处以公共云生产环境为例）
	OperationTimeout    time.Duration     `json:"operation_timeout" mapstructure:"operation_timeout"`   // AccessKey 阿里云身份验证，在阿里云服务器管理控制台创建
	ConnectionTimeout   time.Duration     `json:"connection_timeout" mapstructure:"connection_timeout"` // SecretKey 阿里云身份验证，在阿里云服务器管理控制台创建
	CustomMetricsLabels map[string]string `json:"custom_metrics_labels" mapstructure:"custom_metrics_labels"`
	Timeout             time.Duration     `json:"timeout" mapstructure:"timeout"`
}
