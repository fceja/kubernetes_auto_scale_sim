package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

type Config struct {
	BrokerAddresses []string
	ConsumerGroupId string
	TopicName       string
}

type MyConsumer struct{}

// Load environment variables into config.
func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		BrokerAddresses: strings.Split(os.Getenv("BROKER_ADDRESSES"), ","),
		ConsumerGroupId: os.Getenv("CONSUMER_GROUP_ID"),
		TopicName:       os.Getenv("TOPIC_NAME"),
	}
}

// Runs before consumer starts processing messages.
func (consumer *MyConsumer) Setup(sess sarama.ConsumerGroupSession) error {
	log.Printf("Partition assigned: %v", sess.Claims())

	return nil
}

// Runs after consumer stops processing messages.
func (consumer *MyConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// Consume messages from assigned partition in consumer group,
// and manually commit offset after processing each message.
func (consumer *MyConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// process message
		processMessage(message)

		// commit offset
		sess.MarkMessage(message, "")
		sess.Commit()
	}

	return nil
}

// Processes message.
func processMessage(msg *sarama.ConsumerMessage) {
	log.Printf("Partition %v, Offset %v", msg.Partition, msg.Offset)
	log.Printf("%s", string(msg.Value))
	time.Sleep(1 * time.Second)
}

// Main
func main() {
	// load config
	config := loadConfig()
	log.Printf("config: %+v", config)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	// create kafka client
	client, err := sarama.NewClient(config.BrokerAddresses, saramaConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka client: %v", err)
	}
	defer client.Close()

	// create consumer group
	consumerGroup, err := sarama.NewConsumerGroupFromClient(config.ConsumerGroupId, client)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// consume messages
	consumer := &MyConsumer{}
	ctx := context.Background()
	for {
		err = consumerGroup.Consume(ctx, []string{config.TopicName}, consumer)
		if err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}
}
