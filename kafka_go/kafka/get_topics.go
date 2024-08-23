package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

func GetTopics(brokers []string) []string {
	// create sarama config
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// create sarama client
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatalf("\n\nError creating kafka client: %v", err)
	}
	defer client.Close()

	// get topics
	topics, err := client.Topics()
	if err != nil {
		log.Fatalf("\n\nError retrieving topics: %v", err)
	}

	return topics
}
