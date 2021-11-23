package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var ServerConfigSettings = &ServerConfig{}

func InitConfig(path string, serverName string) {
	flag.Parse()
	flag.StringVar(&ServerConfigSettings.Server.ServerHost, "host", "0.0.0.0", "host 需要配置")
	flag.IntVar(&ServerConfigSettings.Server.ServerPort, "port", 8080, "port 需要配置")
	flag.StringVar(&ServerConfigSettings.Server.ServerName, "server_name", "", "server_name 需要配置")
	flag.StringVar(&ServerConfigSettings.Server.Path, "path", "", "path 需要配置")
	if ServerConfigSettings.Server.ServerName == "" {
		ServerConfigSettings.Server.ServerName = serverName
	}
	if ServerConfigSettings.Server.Path == "" {
		ServerConfigSettings.Server.Path = path
	}
	initServerConfigSettings()
}

func initServerConfigSettings() {
	config := viper.New()
	config.AddConfigPath(ServerConfigSettings.Server.Path)
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := config.Unmarshal(ServerConfigSettings); err != nil {
		panic(err)
	}
	config.WatchConfig()
	config.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件已经修改了 %s", in.String())
		if err := config.Unmarshal(ServerConfigSettings); err != nil {
			panic(err)
		}
	})

}

func (s *ServerConfig) InitServerConfigSettings() {
	config := viper.New()
	ServerConfigSettings = s
	config.AddConfigPath(s.Server.Path)
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := config.Unmarshal(s); err != nil {
		panic(err)
	}
	config.WatchConfig()
	config.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件已经修改了 %s", in.String())
		if err := config.Unmarshal(s); err != nil {
			panic(err)
		}
	})

}
