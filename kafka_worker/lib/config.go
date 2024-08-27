package lib

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	AppEnv          string
	BrokerAddresses []string
	ConsumerGroupId string
	LogFilePath     string
	LogLevel        string
	TopicName       string
}

// Load environment variables into config.
func LoadConfig() Config {
	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// verify env vars exist
	if os.Getenv("APP_ENV") != "development" {
		zap.L().Fatal("Environment variable must exist.")
	}
	if os.Getenv("LOG_FILE_PATH") == "" {
		zap.L().Fatal("Environment variable must exist.")
	}
	if os.Getenv("LOG_LEVEL") == "" {
		zap.L().Fatal("Environment variable must exist.")
	}

	return Config{
		AppEnv:          os.Getenv("APP_ENV"),
		BrokerAddresses: strings.Split(os.Getenv("BROKER_ADDRESSES"), ","),
		ConsumerGroupId: os.Getenv("CONSUMER_GROUP_ID"),
		LogFilePath:     os.Getenv("LOG_FILE_PATH"),
		LogLevel:        os.Getenv("LOG_LEVEL"),
		TopicName:       os.Getenv("TOPIC_NAME"),
	}
}
