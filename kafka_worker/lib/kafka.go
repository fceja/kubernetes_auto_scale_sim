package lib

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type MyConsumer struct{}

// Create kafka client.
func CreateKafkaClient(config Config) sarama.Client {
	// create sarama config
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	// create kafka client
	client, err := sarama.NewClient(config.BrokerAddresses, saramaConfig)
	if err != nil {
		zap.L().Fatal("Error creating Kafka client.", zap.Error(err))
	}

	return client
}

// Create consumer group.
func CreateConsumerGroup(consumerGroupId string, client sarama.Client) sarama.ConsumerGroup {
	consumerGroup, err := sarama.NewConsumerGroupFromClient(consumerGroupId, client)
	if err != nil {
		zap.L().Fatal("Failed creating consumer group.", zap.Error(err))
	}

	return consumerGroup
}

// Runs before consumer starts processing messages.
func (consumer *MyConsumer) Setup(sess sarama.ConsumerGroupSession) error {
	zap.L().Debug("Running consumer setup.", zap.String("partitionsAssigned", fmt.Sprintf("%+v", sess.Claims())))

	return nil
}

// Runs after consumer stops processing messages.
func (consumer *MyConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	zap.L().Debug("Running consumer cleanup.")

	return nil
}

// Consume messages from assigned partition in consumer group,
// and manually commit offset after processing each message.
func (consumer *MyConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// process message
		processMessage(message)

		// commit offset
		sess.MarkMessage(message, "")
		sess.Commit()
	}

	return nil
}

// Process message.
func processMessage(msg *sarama.ConsumerMessage) {
	zap.L().Debug("Processing claims message.",
		zap.String("sleeping", fmt.Sprintf("%v", 5*time.Second)),
		zap.String("partition", fmt.Sprintf("%+v", msg.Partition)),
		zap.String("offset", fmt.Sprintf("%+v", msg.Offset)),
		zap.String("jsonEncoded", string(msg.Value)),
	)

	// sleep
	time.Sleep(3 * time.Second)
}
