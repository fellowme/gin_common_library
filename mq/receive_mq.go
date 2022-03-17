package mq

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	gin_util "github.com/fellowme/gin_common_library/util"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

// ReceivePulsarMqMixMessage 出现  执行func  无序 异步执行
func ReceivePulsarMqMixMessage(pulsarOptions pulsar.ConsumerOptions, f func(message pulsar.Message), stopChan chan error, pool *ants.Pool) {
	consumer, err := getPulsarClient().Subscribe(pulsarOptions)
	if err != nil {
		zap.L().Error("ReceivePulsarMqMixMessage pulsarClient Subscribe error ", zap.Any("error", err))
		stopChan <- err
		return
	}
	defer consumer.Close()
	for {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			zap.L().Error("ReceivePulsarMqMixMessage pulsarClient Receive error ", zap.Any("error", err))
			stopChan <- err
			break
		}
		if msg != nil {
			consumer.Ack(msg)
			err := pool.Submit(func() {
				zap.L().Info("ReceivePulsarMqMixMessage pulsarClient  pool Submit function", zap.Any("msg", string(msg.Payload())))
				f(msg)
			})
			if err != nil {
				zap.L().Error("ReceivePulsarMqMixMessage pulsarClient pool Submit error ", zap.Any("error", err))
				stopChan <- err
				break
			}
		}
	}
}

// ReceivePulsarMqMessage 出现 func 顺序执行 同步执行
func ReceivePulsarMqMessage(pulsarOptions pulsar.ConsumerOptions, f func(message pulsar.Message), stopChan chan error) {
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
			zap.L().Info("ReceivePulsarMqMessage pulsarClient Receive execute function start ", zap.Any("msg", msg), zap.String("start_time", gin_util.NowTimeToString()))
			f(msg)
			zap.L().Info("ReceivePulsarMqMessage pulsarClient Receive execute function end  ", zap.Any("msg", msg), zap.String("end_time", gin_util.NowTimeToString()))
		}
	}
}
