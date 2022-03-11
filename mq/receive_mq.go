package mq

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	gin_config "github.com/fellowme/gin_common_library/config"
	gin_const "github.com/fellowme/gin_common_library/const"
	"go.uber.org/zap"
	"time"
)

func ReceivePulsarMqMessage(pulsarOptions pulsar.ConsumerOptions, f func(message pulsar.Message), stopChan chan error) {
	consumer, err := getPulsarClient().Subscribe(pulsarOptions)
	if err != nil {
		zap.L().Error("ReceivePulsarMqMessage pulsarClient Subscribe error ", zap.Any("error", err))
		stopChan <- err
		return
	}
	defer consumer.Close()
	ctx, cancel := context.WithTimeout(context.Background(), gin_config.ServerConfigSettings.PulsarMqConf.Timeout*time.Second)
	defer cancel()
	for {
		msg, err := consumer.Receive(ctx)
		if err != nil {
			zap.L().Error("ReceivePulsarMqMessage pulsarClient Receive error ", zap.Any("error", err))
			stopChan <- err
			break
		}
		if msg != nil {
			consumer.Ack(msg)
			startTime := time.Now()
			zap.L().Info("ReceivePulsarMqMessage pulsarClient start  ......... ", zap.Any("msg", msg),
				zap.Any("start_time", startTime.Format(gin_const.TimeFormat)))
			f(msg)
			costTime := time.Now().Sub(startTime).Seconds()
			zap.L().Info("ReceivePulsarMqMessage execute success  end ........", zap.Any("cost_time", costTime),
				zap.Any("end_time", time.Now().Format(gin_const.TimeFormat)), zap.Any("msg", msg))
		}
	}
}
