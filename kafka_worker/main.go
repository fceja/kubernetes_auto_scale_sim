package main

import (
	"context"
	"fmt"
	"kafka_worker/lib"

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

	// create kafka client
	client := lib.InitKafkaClient(config)
	defer client.Close()

	// create consumer group
	consumerGroup := lib.CreateConsumerGroup(config.ConsumerGroupId, client)
	defer consumerGroup.Close()

	// consume messages
	consumer := &lib.MyConsumer{}
	ctx := context.Background()
	for {
		err := consumerGroup.Consume(ctx, []string{config.TopicName}, consumer)
		if err != nil {
			zap.L().Fatal("Error consuming messages.", zap.Error(err))
		}
	}
}
