package main

import (
	"fmt"
	"kafka_create_topic/lib"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

func main() {
	// load config
	config, err := lib.LoadConfig()
	if err != nil {
		zap.L().Fatal("Error loading loading config.", zap.Error(err))
	}

	// create zapLogger
	zapLogger, err := lib.SetupLogger(config)
	if err != nil {
		zap.L().Fatal("Error setting up logger.", zap.Error(err))
	}
	zap.ReplaceGlobals(zapLogger) // replace global logger with zap logger
	defer zapLogger.Sync()

	// log config settings
	zap.L().Debug("Configuration values set.", zap.String("config", fmt.Sprintf("%+v", config)))

	// init sarama config
	saramaConfig := sarama.NewConfig()
	client, err := lib.CreateKafkaClient(config.BrokerAddresses, saramaConfig)
	if err != nil {
		zap.L().Fatal("Error creating Kafka client.", zap.Error(err))
	}

	// if topic does not exist, create
	exists := lib.TopicExists(client, config.TopicName)
	if !exists {
		lib.CreateTopic(config, client, config.TopicName)
	} else {
		zap.L().Info("Topic already exists.", zap.String("topicName", config.TopicName))
	}
}
