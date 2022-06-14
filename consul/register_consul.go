package consul

import (
	"errors"
	"fmt"
	gin_config "github.com/fellowme/gin_common_library/config"
	consul "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"time"
)

type ServiceConsul struct {
	Id       string
	Name     string
	Port     int
	Address  string
	IsSecure bool
}

func RegisterGrpcConsul(serviceStruct ServiceConsul) error {
	if clientAgent == nil {
		return errors.New("consul yaml error")
	}
	err := clientAgent.ServiceRegister(&consul.AgentServiceRegistration{
		ID:   serviceStruct.Id,
		Name: serviceStruct.Name,
		Tags: []string{
			"grpc",
		},
		Port:    serviceStruct.Port,
		Address: serviceStruct.Address,
		Check: &consul.AgentServiceCheck{
			Interval: (gin_config.ServerConfigSettings.ConsulConfig.IntervalTime * time.Second).String(),
			GRPC:     fmt.Sprintf("%s:%d/%s", serviceStruct.Address, serviceStruct.Port, serviceStruct.Name),
			Timeout:  (gin_config.ServerConfigSettings.ConsulConfig.TimeOut * time.Second).String(),
		},
	})
	return err
}

func RegisterWebConsul(serviceStruct ServiceConsul) error {
	if clientAgent == nil {
		zap.L().Error("RegisterWebConsul consul yaml error ")
		return errors.New("consul yaml error")
	}
	schema := "http"
	if serviceStruct.IsSecure {
		schema = "https"
	}
	err := clientAgent.ServiceRegister(&consul.AgentServiceRegistration{
		ID:   serviceStruct.Id,
		Name: serviceStruct.Name,
		Tags: []string{
			"http",
		},
		Port:    serviceStruct.Port,
		Address: serviceStruct.Address,
		Check: &consul.AgentServiceCheck{
			Interval: (gin_config.ServerConfigSettings.ConsulConfig.IntervalTime * time.Second).String(),
			HTTP:     fmt.Sprintf("%s://%s:%d/%s", schema, serviceStruct.Address, serviceStruct.Port, "api/v1/ping"),
			Timeout:  (gin_config.ServerConfigSettings.ConsulConfig.TimeOut * time.Second).String(),
			Method:   "HEAD",
		},
	})
	return err
}

func UnRegisterConsul(serviceId string) {
	if clientAgent == nil {
		zap.L().Error("UnRegisterConsul error", zap.Any("error", "clientAgent is nil"))
		return
	}
	err := clientAgent.ServiceDeregister(serviceId)
	if err != nil {
		zap.L().Error("UnRegisterConsul error", zap.Any("error", err))
	}
	return
}
