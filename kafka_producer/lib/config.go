package lib

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	AppEnv          string
	BrokerAddresses []string
	LogFilePath     string
	LogLevel        string
	MessageLimit    int
	SleepTimeout    time.Duration
	TopicName       string
}

// Load environment vars and into config.
// Apply conversions, if needed
func LoadConfig() Config {
	// load env vars
	err := godotenv.Load()
	if err != nil {
		zap.L().Fatal("Error loading .env file")
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

	// conversions
	messageLimit := ConvertStrToInt(os.Getenv("MESSAGE_LIMIT"))
	sleepTimeout := ConvertStrToInt(os.Getenv("SLEEP_TIMEOUT"))

	return Config{
		AppEnv:          os.Getenv("APP_ENV"),
		BrokerAddresses: strings.Split(os.Getenv("BROKER_ADDRESSES"), ","),
		LogLevel:        os.Getenv("LOG_LEVEL"),
		LogFilePath:     os.Getenv("LOG_FILE_PATH"),
		MessageLimit:    messageLimit,
		SleepTimeout:    time.Duration(sleepTimeout) * time.Millisecond,
		TopicName:       os.Getenv("TOPIC_NAME"),
	}
}
