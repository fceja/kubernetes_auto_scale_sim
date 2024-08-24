package main

import (
	"fmt"
	"kafka_go/kafka"
	"log"
	"time"
)

func main() {
	// topic to create or retrieve
	var topicName string = "example_topic_2"

	// define kafka broker (kafka docker container)
	brokers := []string{"localhost:9092"}

	// check if topic exists
	exists := kafka.TopicExists(brokers, topicName)
	fmt.Printf("Topic exists: %v\n", exists)

	// if topic does not exist, create
	if !exists {
		fmt.Print("Creating topic.")
		kafka.CreateTopic(brokers, topicName)
	}

	// add message to topic
	sendMessage := kafka.Message{
		ID:        123,
		Message:   "This is the message to add, again.",
		Name:      "John D",
		Timestamp: time.Now(),
	}
	kafka.AddMessageToTopic(brokers, topicName, sendMessage)

	// read messages from topic
	messages, err := kafka.ReadMessagesFromTopic(brokers, topicName)
	if err != nil {
		log.Fatalf("Error retrieving messages: %v", err)
	}

	fmt.Printf("\nTopic Messages for %v:", topicName)
	for _, msg := range messages {
		fmt.Printf("\nID=%d, Message=%s, Name=%s, Timestamp=%s", msg.ID, msg.Message, msg.Name, msg.Timestamp)
	}
}
