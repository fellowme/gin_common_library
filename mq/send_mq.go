package mq

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	gin_config "github.com/fellowme/gin_common_library/config"
	"go.uber.org/zap"
)

func SendPulsarMqMessage(pulsarOptions pulsar.ProducerOptions, message pulsar.ProducerMessage) (pulsar.MessageID, error) {
	producer, err := getPulsarClient().CreateProducer(pulsarOptions)
	if err != nil {
		zap.L().Error("SendPulsarMq pulsarClient CreateProducer error ", zap.Any("error", err))
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), gin_config.ServerConfigSettings.PulsarMqConf.Timeout)
	defer cancel()
	zap.L().Info("SendPulsarMqMessage start ", zap.Any("pulsarOptions", pulsarOptions), zap.Any("message", message))
	messageId, err := producer.Send(ctx, &message)
	defer producer.Close()
	zap.L().Info("SendPulsarMqMessage end ", zap.Any("pulsarOptions", pulsarOptions), zap.Any("message", message), zap.Any("messageId", messageId))
	return messageId, err
}
