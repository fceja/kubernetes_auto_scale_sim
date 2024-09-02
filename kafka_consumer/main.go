package main

import (
	"context"
	"fmt"
	"kafka_consumer/lib"
	"log"
	"time"

	"go.uber.org/zap"
)

var WAIT_TIME time.Duration = 5 * time.Second

func main() {
	// load config
	config := lib.LoadConfig()

	// create zapLogger
	zapLogger, err := lib.SetupZapLogger(config)
	if err != nil {
		log.Fatal("Error setting up logger.")
	}
	zap.ReplaceGlobals(zapLogger) // replace global logger with zap logger
	defer zapLogger.Sync()

	// log config settings
	zap.L().Debug("Configuration values set.", zap.String("config", fmt.Sprintf("%+v", config)))

	// create kafka client
	client := lib.InitKafkaClient(config)
	defer client.Close()

	// if topic does not exist, create
	var limit int8 = 4
	var retryWait time.Duration = 10 * time.Second

	for i := int8(0); i < limit; i++ {
		exists := lib.CheckIfTopicExists(client, config.TopicName)
		if !exists {
			// topic does not exist, wait before checking again.
			zap.L().Warn(fmt.Sprintf("Trying again in '%v' seconds. '%v' attempt(s) left.", retryWait.String(), limit-i))
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
