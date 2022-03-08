package mq

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	gin_config "github.com/fellowme/gin_common_library/config"
	"go.uber.org/zap"
)

func SendPulsarMqMessage(pulsarOptions pulsar.ProducerOptions, message pulsar.ProducerMessage) error {
	producer, err := getPulsarClient().CreateProducer(pulsarOptions)
	if err != nil {
		zap.L().Error("SendPulsarMq pulsarClient CreateProducer error ", zap.Any("error", err))
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), gin_config.ServerConfigSettings.PulsarMqConf.Timeout)
	defer cancel()
	_, err = producer.Send(ctx, &message)
	defer producer.Close()
	return err
}

func ReceivePulsarMqMessage(pulsarOptions pulsar.ConsumerOptions, f func(message pulsar.Message)) error {
	consumer, err := getPulsarClient().Subscribe(pulsarOptions)
	if err != nil {
		zap.L().Error("ReceivePulsarMqMessage pulsarClient Subscribe error ", zap.Any("error", err))
		return err
	}
	defer consumer.Close()
	ctx, cancel := context.WithTimeout(context.Background(), gin_config.ServerConfigSettings.PulsarMqConf.Timeout)
	defer cancel()
	for {
		msg, err := consumer.Receive(ctx)
		if err != nil {
			zap.L().Error("ReceivePulsarMqMessage pulsarClient Receive error ", zap.Any("error", err))
		}
		if msg != nil {
			consumer.Ack(msg)
			zap.L().Info("ReceivePulsarMqMessage pulsarClient Receive info ", zap.Any("msg", msg))
			f(msg)
		}
	}
}
