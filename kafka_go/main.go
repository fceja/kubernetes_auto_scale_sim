package main

import (
	"fmt"
	"kafka_go/kafka"
)

func main() {
	// topic to create or retrieve
	var topicName string = "example_topic_2"

	// define kafka broker (kafka docker container)
	brokers := []string{"localhost:9092"}

	// check if topic exists
	exists := kafka.TopicExists(brokers, topicName)
	fmt.Printf("Topic exists: %v", exists)

	// if topic does not exist, create
	if !exists {
		fmt.Print("Creating topic.")
		kafka.CreateTopic(brokers, topicName)
	}

}
