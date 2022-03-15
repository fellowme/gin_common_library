package mq

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

func ReceivePulsarMqMessage(pulsarOptions pulsar.ConsumerOptions, f func(message pulsar.Message), stopChan chan error, pool *ants.Pool) {
	consumer, err := getPulsarClient().Subscribe(pulsarOptions)
	if err != nil {
		zap.L().Error("ReceivePulsarMqMessage pulsarClient Subscribe error ", zap.Any("error", err))
		stopChan <- err
		return
	}
	defer consumer.Close()
	for {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			zap.L().Error("ReceivePulsarMqMessage pulsarClient Receive error ", zap.Any("error", err))
			stopChan <- err
			break
		}
		if msg != nil {
			consumer.Ack(msg)
			err := pool.Submit(func() {
				zap.L().Info("ReceivePulsarMqMessage pulsarClient  pool Submit function", zap.Any("msg", string(msg.Payload())))
				f(msg)
			})
			if err != nil {
				zap.L().Error("ReceivePulsarMqMessage pulsarClient pool Submit error ", zap.Any("error", err))
				stopChan <- err
				break
			}
		}
	}
}
