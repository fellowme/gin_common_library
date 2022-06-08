package mq

import (
	"github.com/apache/pulsar-client-go/pulsar"
	gin_config "github.com/fellowme/gin_common_library/config"
	"go.uber.org/zap"
	"time"
)

var pulsarClient pulsar.Client

func InitPulsarClient() {
	var err error
	if gin_config.ServerConfigSettings.PulsarMqConf.PulsarUrl != "" {
		pulsarClient, err = pulsar.NewClient(pulsar.ClientOptions{
			URL:               gin_config.ServerConfigSettings.PulsarMqConf.PulsarUrl,
			OperationTimeout:  gin_config.ServerConfigSettings.PulsarMqConf.OperationTimeout * time.Second,
			ConnectionTimeout: gin_config.ServerConfigSettings.PulsarMqConf.ConnectionTimeout * time.Second,
		})
		if err != nil {
			zap.L().Error("Could not instantiate Pulsar client", zap.Any("error", err))
		}
	}

}

func ClosePulsarClient() {
	if pulsarClient != nil {
		pulsarClient.Close()
	}

}

func getPulsarClient() pulsar.Client {
	if pulsarClient == nil {
		InitPulsarClient()
	}
	return pulsarClient
}
