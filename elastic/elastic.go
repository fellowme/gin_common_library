package elastic

import (
	"context"
	"fmt"
	gin_config "github.com/fellowme/gin_commom_library/config"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

var elasticClient *elastic.Client

func InitElastic() {
	var err error
	url := fmt.Sprintf("http://%s:%d", gin_config.ServerConfigSettings.ElasticConfig.Host, gin_config.ServerConfigSettings.ElasticConfig.Port)
	elasticClient, err = elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(url),
		// elastic.SetBasicAuth(gin_config.ServerConfigSettings.ElasticConfig.UserName, gin_config.ServerConfigSettings.ElasticConfig.Password),
	)
	if err != nil {
		zap.L().Error("es init error", zap.Any("error", err), zap.String("url", url))
	}
	_, _, pingError := elasticClient.Ping(url).Do(context.Background())
	if pingError != nil {
		zap.L().Error("es ping error", zap.Any("error", err))
	}
}

func GetElasticClient() *elastic.Client {
	if elasticClient == nil {
		InitElastic()
	}
	return elasticClient
}
