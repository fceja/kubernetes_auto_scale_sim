package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Message struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func AddMessageToTopic(brokers []string, topicName string, message Message) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // ensure producer waits for acknowledgment

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}
	defer producer.Close()

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(jsonMessage),
	})
	if err != nil {
		log.Fatalf("Failed to send message to kafka: %v", err)
	}

	return nil
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

func ReadMessagesFromTopic(brokers []string, topicName string) ([]Message, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	partition := int32(0)         // adjust if topic has multiple partitions
	offset := sarama.OffsetOldest // start from the beginning (oldest)

	partitionConsumer, err := consumer.ConsumePartition(topicName, partition, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to consume partition: %w", err)
	}
	defer partitionConsumer.Close()

	var messages []Message
	timeout := time.After(2 * time.Second) // Adjust the timeout as needed

	for {
		select {
		case message, ok := <-partitionConsumer.Messages():
			if !ok {
				// channel closed
				return messages, nil
			}

			var msg Message
			err := json.Unmarshal(message.Value, &msg)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			messages = append(messages, msg)
		case <-timeout:
			// timeout reached
			log.Println("Timeout reached, stopping message retrieval.")
			return messages, nil
		}
	}
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
