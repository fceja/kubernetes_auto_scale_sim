package main

import (
	"fmt"
	"kafka_producer/lib"
	"time"

	"go.uber.org/zap"
)

func main() {
	// load config
	config := lib.LoadConfig()

	// create zapLogger
	zapLogger, err := lib.SetupLogger(config)
	if err != nil {
		zap.L().Fatal("Error setting up logger.", zap.Error(err))
	}
	zap.ReplaceGlobals(zapLogger) // replace global logger with zap logger
	defer zapLogger.Sync()

	// log config settings
	zap.L().Debug("Configuration values set.", zap.String("config", fmt.Sprintf("%+v", config)))

	// init kafka producer client
	client, err := lib.InitKafkaProducerClient(config.BrokerAddresses)
	if err != nil {
		zap.L().Fatal("Error creating Kafka client.", zap.Error(err))
	}

	// create and add messages to topic until limit reached
	for i := 0; i < config.MessageLimit; i++ {
		sendMessage := lib.Message{
			ID:        i,
			Message:   fmt.Sprintf("This is the message %v.", i),
			Name:      "John D",
			Timestamp: time.Now(),
		}

		// send message to topic
		lib.AddMessageToTopic(client, sendMessage, config.TopicName)
		zap.L().Info("Message added to topic.")

		// sleep
		zap.L().Info(fmt.Sprintf("Sleep for %v.", config.SleepTimeout))
		time.Sleep(config.SleepTimeout)
	}
}
