package lib

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// Checks if slice contains item
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

func CreateKafkaClient(brokerAddresses []string, saramaConfig *sarama.Config) (sarama.Client, error) {
	zap.L().Info("Creating Kafka client.")

	saramaClient, err := sarama.NewClient(brokerAddresses, saramaConfig)
	if err != nil {
		zap.L().Fatal("Error creating kafka client.", zap.Error(err))
	}

	zap.L().Info("Kafka client created successfully.")

	return saramaClient, nil
}

// Creates kafka topic
func CreateTopic(config Config, client sarama.Client, topicName string) {
	zap.L().Info("Creating topic.")

	// create sarama admin client
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		zap.L().Fatal("Error creating kafka admin client.", zap.Error(err))
	}
	defer admin.Close()

	// define the topic details
	topicDetail := sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	// create topic
	err = admin.CreateTopic(topicName, &topicDetail, false)
	if err != nil {
		zap.L().Fatal("Error creating topic.", zap.Error(err))

	}

	zap.L().Info("Topic created successfully.", zap.String("topicName", topicName))
}

// Check if kafka topic already exists
func TopicExists(client sarama.Client, topicName string) bool {
	topics, err := client.Topics()
	if err != nil {
		zap.L().Fatal("Error retrieving topics.", zap.Error(err))
	}

	var result bool = contains(topics, topicName)
	zap.L().Info("Check if topic exists already.", zap.Bool("exists", result))

	return result
}
