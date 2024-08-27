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

// load environment vars and apply conversions if needed
func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		zap.L().Fatal("Error loading .env file")
	}

	var appEnv string = os.Getenv("APP_ENV")
	if appEnv != "development" {
		zap.L().Fatal("Environment variable must exist.")
	}

	var logFilePath string = os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		zap.L().Fatal("Environment variable must exist.")
	}

	var logLevel string = os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		zap.L().Fatal("Environment variable must exist.")
	}

	messageLimit := ConvertStrToInt(os.Getenv("MESSAGE_LIMIT"))
	sleepTimeout := ConvertStrToInt(os.Getenv("SLEEP_TIMEOUT"))

	return Config{
		AppEnv:          appEnv,
		BrokerAddresses: strings.Split(os.Getenv("BROKER_ADDRESSES"), ","),
		LogLevel:        logLevel,
		LogFilePath:     logFilePath,
		MessageLimit:    messageLimit,
		SleepTimeout: time.Duration(sleepTimeout) * time.Millisecond,
		TopicName:    os.Getenv("TOPIC_NAME"),
	}
}
