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
	zapLogger, err := lib.SetupZapLogger(config)
	if err != nil {
		zap.L().Fatal("Error setting up logger.", zap.Error(err))
	}
	zap.ReplaceGlobals(zapLogger) // replace global logger with zap logger
	defer zapLogger.Sync()

	// log config settings
	zap.L().Debug("Configuration values set.", zap.String("config", fmt.Sprintf("%+v", config)))

	// init kafka producer producerClient
	producerClient, err := lib.InitKafkaProducerClient(config.BrokerAddresses)
	if err != nil {
		zap.L().Fatal("Error creating Kafka client.", zap.Error(err))
	}

	// create and add messages to topic until retryLimit reached
	var attemptCount = 0
	var limit = 2
	var retryWait time.Duration = 5 * time.Second

	for i := 0; i < config.MessageLimit; i++ {
		sendMessage := lib.Message{
			ID:        i,
			Message:   fmt.Sprintf("This is the message %v.", i),
			Name:      "John D",
			Timestamp: time.Now(),
		}

		// send message to topic
		err := lib.AddMessageToTopic(producerClient, sendMessage, config.TopicName)
		if err != nil {
			if attemptCount < limit {
				zap.L().Warn(fmt.Sprintf("Re-attempting in %v seconds. %v attempt(s) left.", retryWait.Seconds(), limit-attemptCount))

				time.Sleep(retryWait)

				attemptCount++
				continue
			}
			zap.L().Error("Retry limit reached, exiting.")
			break
		}

		zap.L().Info("Message added to topic.")
		attemptCount = 0 // reset

		// sleep
		zap.L().Info(fmt.Sprintf("Sleep for %v.", config.SleepTimeout))
		time.Sleep(config.SleepTimeout)
	}
}
