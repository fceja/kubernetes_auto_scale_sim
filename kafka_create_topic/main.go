package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

type Config struct {
	BrokerAddresses        []string
	TopicName              string
	TopicNumberPartitions  int32
	TopicReplicationFactor int16
}

// Section: helper funcs
// Checks if slices contain item.
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

// Converts string to int.
func convertStrToInt(inputStr string) int {
	inputInt, err := strconv.Atoi(inputStr)
	if err != nil {
		log.Fatalf("Error converting string to int: %v", err)
	}

	return inputInt
}

// Section: kafka funcs
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
func createTopic(config Config, client sarama.Client, topicName string) {
	log.Println("Creating topic.")

	// create sarama admin client
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		log.Fatalf("Error creating kafka admin client: %v", err)
	}
	defer admin.Close()

	// define the topic details
	topicDetail := sarama.TopicDetail{
		NumPartitions:     config.TopicNumberPartitions,
		ReplicationFactor: config.TopicReplicationFactor,
	}

	// create topic
	err = admin.CreateTopic(topicName, &topicDetail, false)
	if err != nil {
		log.Fatalf("Error creating topic: %v", err)
	}
	log.Printf("Topic '%v' created successfully.", topicName)
}

// Check if kafka topic already exists.
func topicExists(client sarama.Client, topicName string) bool {
	topics, err := client.Topics()
	if err != nil {
		log.Fatalf("Error retrieving topics: %v", err)
	}

	return contains(topics, topicName)
}

// Section: init
// Load environment variables to config.
func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	topicNumberPartitions := convertStrToInt(os.Getenv("TOPIC_NUMBER_PARTITIONS"))
	topicReplicationFactor := convertStrToInt(os.Getenv("TOPIC_REPLICATION_FACTOR"))

	var int32topicNumberPartitions = int32(topicNumberPartitions)
	var int16topicNumberPartitions = int16(topicReplicationFactor)

	return Config{
		BrokerAddresses:        strings.Split(os.Getenv("BROKER_ADDRESSES"), ","),
		TopicName:              os.Getenv("TOPIC_NAME"),
		TopicNumberPartitions:  int32topicNumberPartitions,
		TopicReplicationFactor: int16topicNumberPartitions,
	}
}

// Section: main
func main() {
	// load config
	config := loadConfig()
	log.Printf("config: %+v", config)

	// init sarama config
	saramaConfig := sarama.NewConfig()
	client := createKafkaClient(config.BrokerAddresses, saramaConfig)

	// if topic does not exist, create
	exists := topicExists(client, config.TopicName)
	if !exists {
		createTopic(config, client, config.TopicName)
	} else {
		log.Printf("Topic '%v' already exists.", config.TopicName)
	}
}
