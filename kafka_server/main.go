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
	config := lib.LoadConfig()

	// create zapLogger
	zapLogger, err := lib.SetupZapLogger(config)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(zapLogger) // replace global logger with zap logger
	defer zapLogger.Sync()

	zap.L().Debug("Configuration values set.", zap.String("config", fmt.Sprintf("%+v", config)))

	// init sarama config
	var saramaConfig *sarama.Config = sarama.NewConfig()
	var client sarama.Client

	// connect to kafka client
	var retryLimit int8 = 2
	var count int8 = 0
	var waitTime time.Duration = 3 * time.Second

	for {
		if count > retryLimit {
			zap.L().Fatal("Failed to connect with Kafka client.")
		}

		var err error
		client, err = lib.CreateKafkaClient(config.BrokerAddresses, saramaConfig)
		if err != nil {
			zap.L().Error("Error creating Kafka client.", zap.Error(err))
			zap.L().Warn(fmt.Sprintf("\nRetrying in '%v'. Attempts left: '%v'", waitTime, retryLimit-1))
			time.Sleep(waitTime)

			count++
			continue
		}
		break
	}

	// create kafka topic
	var topicName string = "example-topic-2"

	// if topic does not exist, create
	exists := lib.TopicExists(client, topicName)
	if !exists {
		lib.CreateTopic(config, client, topicName)
	} else {
		zap.L().Info("Topic already exists.", zap.String("topicName", topicName))
	}
}
