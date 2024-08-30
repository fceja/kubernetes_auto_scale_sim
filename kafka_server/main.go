package main

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"main.go/lib"
)

func main() {
	// load config
	fmt.Print("Loading config.")
	config := lib.LoadConfig()

	// fmt.Print("Kafka server initialization.")
	// create zapLogger
	zapLogger, err := lib.SetupZapLogger(config)
	if err != nil {
		zap.L().Fatal("Error setting up logger.", zap.Error(err))
	}
	zap.ReplaceGlobals(zapLogger) // replace global logger with zap logger
	defer zapLogger.Sync()

	// log config settings
	zap.L().Debug("Configuration values set.", zap.String("config", fmt.Sprintf("%+v", config)))

	// init sarama config

	var retryLimit int8 = 2
	var count int8 = 0
	var waitTime time.Duration = 3 * time.Second

	var saramaConfig *sarama.Config = sarama.NewConfig()
	var client sarama.Client

	// for i := 0; i < retryLimit; i++ {
	for {
		if count > retryLimit {
			zap.L().Fatal("Failed to connect with Kafka client.")
		}

		var err error
		client, err = lib.CreateKafkaClient(config.BrokerAddresses, saramaConfig)
		if err != nil {
			zap.L().Error("Error creating Kafka client.", zap.Error(err))
			zap.L().Warn(fmt.Sprintf("Retrying in '%v'. Attempts left: '%v'", waitTime, retryLimit-1))
			time.Sleep(waitTime)

			count++
			continue
		}
		break
	}

	var topicName string = "example-topic-1"
	// if topic does not exist, create
	exists := lib.TopicExists(client, topicName)
	if !exists {
		lib.CreateTopic(config, client, topicName)
	} else {
		zap.L().Info("Topic already exists.", zap.String("topicName", topicName))
	}
	// exists := lib.TopicExists(client, config.TopicName)
	// if !exists {
	// 	lib.CreateTopic(config, client, config.TopicName)
	// } else {
	// 	zap.L().Info("Topic already exists.", zap.String("topicName", config.TopicName))
	// }
}
