package consul

import (
	gin_config "github.com/fellowme/gin_common_library/config"
	consul "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

var clientAgent *consul.Agent

func InitConsulClient() {
	if gin_config.ServerConfigSettings.ConsulConfig.ConsulAddress == "" {
		return
	}
	consulConfig := consul.DefaultConfig()
	consulConfig.Address = gin_config.ServerConfigSettings.ConsulConfig.ConsulAddress
	client, err := consul.NewClient(consulConfig)
	if err != nil {
		zap.L().Error("InitConsulClient error", zap.Any("error", err))
	}
	clientAgent = client.Agent()
	return
}
