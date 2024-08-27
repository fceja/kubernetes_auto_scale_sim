package lib

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

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
func AddMessageToTopic(client *KafkaProducerClient, message Message, topicName string) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, _, err = client.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(jsonMessage),
	})
	if err != nil {
		zap.L().Fatal("Failed to send message to kafka.", zap.Error(err))
	}

	return nil
}

// init kafka producer client
func InitKafkaProducerClient(brokerAddresses []string) (*KafkaProducerClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // ensure producer waits for acknowledgment

	producer, err := sarama.NewSyncProducer(brokerAddresses, config)
	if err != nil {
		zap.L().Fatal("Failed to start producer.", zap.Error(err))
	}

	return &KafkaProducerClient{
		Config:   config,
		Producer: producer,
	}, nil
}
