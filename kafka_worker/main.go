package main

import (
	"context"
	"fmt"
	"kafka_worker/lib"
	"log"
	"time"

	"go.uber.org/zap"
)

var WAIT_TIME time.Duration = 5 * time.Second

func main() {
	// load config
	config := lib.LoadConfig()

	// create zapLogger
	// fmt.Println("\nCreating logger")
	zapLogger, err := lib.SetupZapLogger(config)
	if err != nil {
		// zap.L().Fatal("Error setting up logger.", zap.Error(err))
		log.Fatal("Error setting up logger.")
	}
	zap.ReplaceGlobals(zapLogger) // replace global logger with zap logger
	defer zapLogger.Sync()
	// zap.L().Debug("Zap logger created")

	// log config settings
	zap.L().Debug("Configuration values set.", zap.String("config", fmt.Sprintf("%+v", config)))

	// create kafka client
	// zap.L().Debug("InitKafkaClient before")
	time.Sleep(10 * time.Second)
	zap.L().Debug("sleeping 20 secs")

	client := lib.InitKafkaClient(config)
	defer client.Close()
	// zap.L().Debug(fmt.Sprintf("client 3 - %v", client))
	// zap.L().Debug(fmt.Sprintf("client 4 %+v", client))
	// zap.L().Debug("InitKafkaClient after")

	// if topic does not exist, create
	var attemptLimit int8 = 4

	for i := int8(0); i < attemptLimit; i++ {
		exists := lib.CheckIfTopicExists(client, config.TopicName)
		if !exists {
			// topic does not exist, wait before checking again.
			zap.L().Warn(fmt.Sprintf("Trying again in '%v' seconds. '%v' attempt(s) left.", WAIT_TIME.String(), attemptLimit-i))
			time.Sleep(WAIT_TIME)
			continue

		} else {
			// topic exists
			break
		}
	}

	// create consumer group
	consumerGroup := lib.CreateConsumerGroup(config.ConsumerGroupId, client)
	defer consumerGroup.Close()

	// defining consumer
	zap.L().Debug("Defining consumer.")
	consumer := &lib.MyConsumer{}
	ctx := context.Background()

	// consume messages
	var errLimit int8 = 4
	var errCount int8 = 0

	for {
		zap.L().Info("Attempting to consume messages from Kafka server topic.")
		err := consumerGroup.Consume(ctx, []string{config.TopicName}, consumer)
		if err != nil {
			zap.L().Error("Error consuming message from topic:", zap.Error(err))

			errCount++

			if errCount < errLimit {
				zap.L().Warn(fmt.Sprintf("Retrying in %v.", WAIT_TIME.String()))
				time.Sleep(WAIT_TIME)

				continue
			}
			zap.L().Error("Retry count reached. Exited.:", zap.Error(err))
		}
	}
}
