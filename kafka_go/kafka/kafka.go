package kafka

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func AddMessageToTopic(brokers []string, topicName string, message string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // ensure producer waits fro acknowledgment

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}
	defer producer.Close()

	producerMessage := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(producerMessage)
	if err != nil {
		log.Fatalf("Failed to send message to kafka: %v", err)
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topicName, partition, offset)

}

func CreateTopic(brokers []string, topicName string) {
	// create sarama config
	config := sarama.NewConfig()

	// create sarama client
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatalf("Error creating kafka client: %v", err)
	}

	// create sarama admin client
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		log.Fatalf("Error creating kafka admin client: %v", err)
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
		log.Fatalf("Error creating topic: %v", err)
	}
	fmt.Println("Topic created successfully.")
}

func GetTopics(brokers []string) []string {
	// create sarama config
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// create sarama client
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatalf("Error creating kafka client: %v", err)
	}
	defer client.Close()

	// get topics
	topics, err := client.Topics()
	if err != nil {
		log.Fatalf("Error retrieving topics: %v", err)
	}

	return topics
}

func TopicExists(brokers []string, topicName string) bool {
	topics := GetTopics(brokers)

	return contains(topics, topicName)
}

func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}
