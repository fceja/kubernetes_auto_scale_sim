package lib

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	AppEnv                 string
	BrokerAddresses        []string
	LogFilePath            string
	LogLevel               string
	TopicName              string
	TopicNumberPartitions  int32
	TopicReplicationFactor int16
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		zap.L().Fatal("Error loading .env file.", zap.Error(err))
	}

	var appEnv string = os.Getenv("APP_ENV")
	if appEnv != "development" {
		zap.L().Fatal("Environment variable must exist.")
	}

	if os.Getenv("BROKER_ADDRESSES") == "" {
		zap.L().Fatal("Environment variable must exist.")
	}
	var brokerAddresses []string = strings.Split(os.Getenv("BROKER_ADDRESSES"), ",")

	var logFilePath string = os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		zap.L().Fatal("Environment variable must exist.")
	}

	var logLevel string = os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		zap.L().Fatal("Environment variable must exist.")
	}

	var topicName string = os.Getenv("TOPIC_NAME")
	if topicName == "" {
		zap.L().Fatal("Environment variable must exist.")
	}

	var strTopicNumberPartitions string = os.Getenv("TOPIC_NUMBER_PARTITIONS")
	if strTopicNumberPartitions == "" {
		zap.L().Fatal("Environment variable must exist.")
	}
	var intTopicNumberPartitions int = ConvertStrToInt(strTopicNumberPartitions)
	var int32TopicNumberPartitions = int32(intTopicNumberPartitions)

	var strTopicReplicationFactor string = os.Getenv("TOPIC_REPLICATION_FACTOR")
	if strTopicReplicationFactor == "" {
		zap.L().Fatal("Environment variable must exist.")
	}
	var intTopicReplicationFactor int = ConvertStrToInt(strTopicReplicationFactor)
	var int16TopicReplicationFactor = int16(intTopicReplicationFactor)

	return Config{
		AppEnv:                 appEnv,
		BrokerAddresses:        brokerAddresses,
		LogLevel:               logLevel,
		LogFilePath:            logFilePath,
		TopicName:              topicName,
		TopicNumberPartitions:  int32TopicNumberPartitions,
		TopicReplicationFactor: int16TopicReplicationFactor,
	}, nil
}
