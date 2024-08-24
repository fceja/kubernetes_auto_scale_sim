package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

type KafkaClient struct {
	BrokerAddresses []string
	Config        *sarama.Config
	Consumer      sarama.Consumer
	TopicName     string
}

type Message struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

// section: kafka funcs
func newKafkaClient(brokerAddresses []string, topicName string) (*KafkaClient, error) {
	// create sarama config
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	// create consumer
	consumer, err := sarama.NewConsumer(brokerAddresses, config)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	return &KafkaClient{
		Consumer:      consumer,
		Config:        config,
		BrokerAddresses: brokerAddresses,
		TopicName:     topicName,
	}, nil

}

func readMessagesFromTopic(client *KafkaClient) ([]Message, error) {
	// consume topic partition
	partition := int32(0)         // adjust if topic has multiple partitions
	offset := sarama.OffsetOldest // start from the beginning (oldest)
	partitionConsumer, err := client.Consumer.ConsumePartition(client.TopicName, partition, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to consume partition: %w", err)
	}
	defer partitionConsumer.Close()

	// retrieve messages from topic partition
	var messages []Message
	timeout := time.After(1 * time.Second) // Adjust the timeout as needed

	for {
		select {
		case message, ok := <-partitionConsumer.Messages():
			if !ok {
				// channel closed
				return messages, nil
			}

			// parse message to json
			var msg Message
			err := json.Unmarshal(message.Value, &msg)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// store message
			messages = append(messages, msg)
		case <-timeout:
			// timeout reached
			log.Println("Timeout reached, stopping message retrieval.")
			return messages, nil
		}
	}
}

// section: worker funcs
func processMessage(messages []Message) {
	log.Print("Messages:")
	for _, msg := range messages {
		log.Printf("ID=%d, Message=%s, Name=%s, Timestamp=%s", msg.ID, msg.Message, msg.Name, msg.Timestamp)
	}

	time.Sleep(1 * time.Second)
}

// section: main
func main() {
	// get env vars
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// define broker (kafka docker container) and topic name
	var brokerAddresses []string
	brokerAddresses = append(brokerAddresses, strings.Split(os.Getenv("BROKER_ADDRESSES"), ",")...)
	topicName := os.Getenv("TOPIC_NAME")

	// create kafka client
	client, err := newKafkaClient(brokerAddresses, topicName)
	if err != nil {
		log.Fatalf("Error creating kafka client: %v", err)
	}
	log.Printf("kafka client: %+v", client)

	// retrieve messages from topic
	messages, err := readMessagesFromTopic(client)
	if err != nil {
		log.Fatalf("Error retrieving messages: %v", err)
	}

	// process messages
	processMessage(messages)

	// close open connections
	defer func() {
		if client.Consumer != nil {
			client.Consumer.Close()
		}
	}()
}
