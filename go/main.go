package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	// define kafka broker address (kafka docker container )
	brokers := []string{"localhost:9092"}
	fmt.Printf("\n\nbrokers: %v", brokers)

	// create sarama config
	config := sarama.NewConfig()
	fmt.Printf("\n\nconfig: %+v", config)

	// create sarama client
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatalf("\n\nError creating kafka client: %v", err)
	}
	fmt.Printf("\n\nclient: %+v", client)

	// create sarama admin client
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		log.Fatalf("\n\nError creating kafka admin client: %v", err)
	}
	defer admin.Close()
	fmt.Printf("\n\nadmin: %+v", admin)

	// define the topic details
	topicName := "example_topic"
	topicDetail := sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	fmt.Printf("\n\ntopicName: %+v", topicName)
	fmt.Printf("\n\ntopicDetails: %+v", topicDetail)

	// create topic
	err = admin.CreateTopic(topicName, &topicDetail, false)
	if err != nil {
		log.Fatalf("\n\nError creating topic: %v", err)
	}
	fmt.Println("\n\nTopic created successfully.")
}
