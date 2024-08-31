package lib

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type MyConsumer struct{}

// Checks if slice contains item
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

// Check if kafka topic already exists
func CheckIfTopicExists(client sarama.Client, topicName string) bool {
	zap.L().Info("Checking if topic exists.")
	topics, err := client.Topics()
	if err != nil {
		zap.L().Fatal("Error retrieving topics.", zap.Error(err))
	}
	zap.L().Debug("Topics available.",
		zap.Strings("topics", topics),
	)

	var exists bool = contains(topics, topicName)

	if exists {
		zap.L().Info(fmt.Sprintf("Topic '%v' exists.", topicName))
	} else {
		zap.L().Error(fmt.Sprintf("Topic '%v' does not exist.", topicName))
	}

	return exists
}

// Create consumer group.
func CreateConsumerGroup(consumerGroupId string, client sarama.Client) sarama.ConsumerGroup {
	zap.L().Debug("Creating consumer group.")
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
		ProcessMessage(message)

		// commit offset
		sess.MarkMessage(message, "")
		sess.Commit()
	}

	return nil
}

func pollForKafkaConnection(brokerAddresses []string, config *sarama.Config) sarama.Client {
	var client sarama.Client

	for {
		var err error
		client, err = sarama.NewClient(brokerAddresses, config)
		if err != nil {
			zap.L().Error("Failed to connect to Kafka.", zap.Error(err))
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	return client
}

// Create kafka client.
func InitKafkaClient(config Config) sarama.Client {
	zap.L().Debug("Creating sarama Kafka client.")

	// create sarama config
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	client := pollForKafkaConnection(config.BrokerAddresses, saramaConfig)

	return client
}
