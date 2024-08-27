package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

type Config struct {
	BrokerAddresses []string
	MessageLimit    int
	SleepTimeout    time.Duration
	TopicName       string
}

type KafkaProducerClient struct {
	Config   *sarama.Config
	Producer sarama.SyncProducer
}

type Message struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

// adds message to topic
func addMessageToTopic(client *KafkaProducerClient, message Message, topicName string) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, _, err = client.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(jsonMessage),
	})
	if err != nil {
		log.Fatalf("Failed to send message to kafka: %v", err)
	}

	return nil
}

// helper
func convertStrToInt(inputStr string) int {
	inputInt, err := strconv.Atoi(inputStr)
	if err != nil {
		log.Fatalf("Error converting string to int: %v", err)
	}

	return inputInt
}

// load environment vars and apply conversions if needed
func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	messageLimit := convertStrToInt(os.Getenv("MESSAGE_LIMIT"))
	sleepTimeout := convertStrToInt(os.Getenv("SLEEP_TIMEOUT"))

	return Config{
		BrokerAddresses: strings.Split(os.Getenv("BROKER_ADDRESSES"), ","),
		MessageLimit:    messageLimit,
		// SleepTimeout:    time.Duration(sleepTimeout) * time.Second,
		SleepTimeout: time.Duration(sleepTimeout) * time.Millisecond,
		TopicName:    os.Getenv("TOPIC_NAME"),
	}
}

// init kafka producer client
func initKafkaProducerClient(brokerAddresses []string) (*KafkaProducerClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // ensure producer waits for acknowledgment

	producer, err := sarama.NewSyncProducer(brokerAddresses, config)
	if err != nil {
		log.Fatalf("Failed to start producer: %v", err)
	}

	return &KafkaProducerClient{
		Config:   config,
		Producer: producer,
	}, nil
}

// main
func main() {
	// load config
	config := loadConfig()
	log.Printf("config %+v", config)

	// init kafka producer client
	client, err := initKafkaProducerClient(config.BrokerAddresses)
	if err != nil {
		log.Fatalf("Error creating kafka client: %v", err)
	}

	// create and add messages to topic until limit reached
	for i := 0; i < config.MessageLimit; i++ {
		// generate message
		sendMessage := Message{
			ID:        i,
			Message:   fmt.Sprintf("This is the message %v", i),
			Name:      "John D",
			Timestamp: time.Now(),
		}

		// send message to topic
		addMessageToTopic(client, sendMessage, config.TopicName)

		// sleep
		log.Printf("Message %v added to topic, sleeping %v.", i, config.SleepTimeout)
		time.Sleep(config.SleepTimeout)
	}
}
