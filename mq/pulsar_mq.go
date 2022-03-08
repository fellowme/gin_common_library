package mq

import (
	"github.com/apache/pulsar-client-go/pulsar"
	gin_config "github.com/fellowme/gin_common_library/config"
	"go.uber.org/zap"
)

var pulsarClient pulsar.Client

func InitPulsarClient() {
	var err error
	pulsarClient, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:               gin_config.ServerConfigSettings.PulsarMqConf.PulsarUrl,
		OperationTimeout:  gin_config.ServerConfigSettings.PulsarMqConf.OperationTimeout,
		ConnectionTimeout: gin_config.ServerConfigSettings.PulsarMqConf.ConnectionTimeout,
	})
	if err != nil {
		zap.L().Error("Could not instantiate Pulsar client", zap.Any("error", err))
	}
}

func ClosePulsarClient() {
	pulsarClient.Close()
}

func getPulsarClient() pulsar.Client {
	if pulsarClient == nil {
		InitPulsarClient()
	}
	return pulsarClient
}
