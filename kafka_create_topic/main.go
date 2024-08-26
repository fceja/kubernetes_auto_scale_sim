package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

type Config struct {
	BrokerAddresses []string
	TopicName       string
}

// helper
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

// Initializes kafka client.
func createKafkaClient(brokerAddresses []string, saramaConfig *sarama.Config) sarama.Client {
	saramaClient, err := sarama.NewClient(brokerAddresses, saramaConfig)
	if err != nil {
		log.Fatalf("Error creating kafka client: %v", err)
	}

	// return saramaClient, nil
	return saramaClient
}

// Creates kafka topic.
func createTopic(client sarama.Client, topicName string) {
	fmt.Println("Creating topic.")

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
	fmt.Printf("Topic '%v' created successfully.", topicName)
}

// Load environment variables to config.
func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		BrokerAddresses: strings.Split(os.Getenv("BROKER_ADDRESSES"), ","),
		TopicName:       os.Getenv("TOPIC_NAME"),
	}
}

// Check if kafka topic already exists.
func topicExists(client sarama.Client, topicName string) bool {
	topics, err := client.Topics()
	if err != nil {
		log.Fatalf("Error retrieving topics: %v", err)
	}

	return contains(topics, topicName)
}

// Main
func main() {
	// load config
	config := loadConfig()
	log.Printf("config: %+v", config)

	saramaConfig := sarama.NewConfig()
	client := createKafkaClient(config.BrokerAddresses, saramaConfig)

	// check if topic exists
	exists := topicExists(client, config.TopicName)
	fmt.Printf("Topic already exists: %v\n", exists)

	// if topic does not exist, create
	if !exists {
		createTopic(client, config.TopicName)
	}
}
