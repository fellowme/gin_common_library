package send_phone

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/fellowme/gin_common_library/config"
	"go.uber.org/zap"
)

func createClient() (*dysmsapi20170525.Client, error) {
	openApiConfig := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: tea.String(config.ServerConfigSettings.AliYunSendPhoneCodeConfig.AccessKeyId),
		// 您的AccessKey Secret
		AccessKeySecret: tea.String(config.ServerConfigSettings.AliYunSendPhoneCodeConfig.AccessKeySecret),
	}
	// 访问的域名
	openApiConfig.Endpoint = tea.String(config.ServerConfigSettings.AliYunSendPhoneCodeConfig.Endpoint)
	result, err := dysmsapi20170525.NewClient(openApiConfig)
	return result, err
}

func AliYunSendSms(sendSmsRequest dysmsapi20170525.SendSmsRequest) error {
	client, err := createClient()
	if err != nil {
		zap.L().Error("AliYunSendSms createClient error", zap.Any("sendSmsRequest", sendSmsRequest), zap.Any("error", err))
		return err
	}
	_, err = client.SendSms(&sendSmsRequest)
	if err != nil {
		zap.L().Error("AliYunSendSms SendSms error", zap.Any("sendSmsRequest", sendSmsRequest), zap.Any("error", err))
	}
	return err
}

func AliYunSendBatchSms(sendBatchSmsRequest dysmsapi20170525.SendBatchSmsRequest) error {
	client, err := createClient()
	if err != nil {
		zap.L().Error("AliYunSendBatchSms createClient error", zap.Any("sendSmsRequest", sendBatchSmsRequest), zap.Any("error", err))
		return err
	}
	_, err = client.SendBatchSms(&sendBatchSmsRequest)
	if err != nil {
		zap.L().Error("AliYunSendBatchSms SendSms error", zap.Any("sendSmsRequest", sendBatchSmsRequest), zap.Any("error", err))
	}
	return err
}
