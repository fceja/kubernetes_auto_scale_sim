package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	AppEnv                 string
	BrokerAddresses        []string
	LogFilePath            string
	TopicName              string
	TopicNumberPartitions  int32
	TopicReplicationFactor int16
}

// Section: helper funcs
// Checks if slices contain item
func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

// Converts string to int
func convertStrToInt(inputStr string) int {
	inputInt, err := strconv.Atoi(inputStr)
	if err != nil {
		zap.L().Fatal("Error converting string to int", zap.Error(err))
	}

	return inputInt
}

// Section: kafka funcs
// Initializes kafka client
func createKafkaClient(brokerAddresses []string, saramaConfig *sarama.Config) sarama.Client {
	zap.L().Info("creating kafka client")

	saramaClient, err := sarama.NewClient(brokerAddresses, saramaConfig)
	if err != nil {
		zap.L().Fatal("Error creating kafka client", zap.Error(err))
	}

	zap.L().Info("kafka client created successfully")

	return saramaClient
}

// Creates kafka topic
func createTopic(config Config, client sarama.Client, topicName string) {
	zap.L().Info("Creating topic.")

	// create sarama admin client
	admin, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		log.Fatalf("Error creating kafka admin client: %v", err)
	}
	defer admin.Close()

	// define the topic details
	topicDetail := sarama.TopicDetail{
		NumPartitions:     config.TopicNumberPartitions,
		ReplicationFactor: config.TopicReplicationFactor,
	}

	// create topic
	err = admin.CreateTopic(topicName, &topicDetail, false)
	if err != nil {
		log.Fatalf("Error creating topic: %v", err)
	}

	zap.L().Info("Topic created successfully.", zap.String("topicName", topicName))

}

// Check if kafka topic already exists
func topicExists(client sarama.Client, topicName string) bool {
	zap.L().Info("checking if topic already exists")

	topics, err := client.Topics()
	if err != nil {
		zap.L().Fatal("Error retrieving topics", zap.Error(err))
	}

	var result bool = contains(topics, topicName)
	zap.L().Info("if topic exists", zap.Bool("exists", result))

	return result
}

// Section: init
// Load environment variables to config.
func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		zap.L().Fatal("Error loading .env file", zap.Error(err))
	}

	if os.Getenv("BROKER_ADDRESSES") == "" {
		zap.L().Fatal("Environment variable must exist.")
	}
	var brokerAddresses []string = strings.Split(os.Getenv("BROKER_ADDRESSES"), ",")

	var logFilePath string = os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		zap.L().Fatal("Environment variable must exist.")
	}

	var logLevel string = os.Getenv("APP_ENV")
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
	var intTopicNumberPartitions int = convertStrToInt(strTopicNumberPartitions)
	var int32TopicNumberPartitions = int32(intTopicNumberPartitions)

	var strTopicReplicationFactor string = os.Getenv("TOPIC_REPLICATION_FACTOR")
	if strTopicReplicationFactor == "" {
		zap.L().Fatal("Environment variable must exist.")
	}
	var intTopicReplicationFactor int = convertStrToInt(strTopicReplicationFactor)
	var int16TopicReplicationFactor = int16(intTopicReplicationFactor)

	return Config{
		AppEnv:                 logLevel,
		BrokerAddresses:        brokerAddresses,
		LogFilePath:            logFilePath,
		TopicName:              topicName,
		TopicNumberPartitions:  int32TopicNumberPartitions,
		TopicReplicationFactor: int16TopicReplicationFactor,
	}
}

func setupLogger(config Config) *zap.Logger {
	// create dir if does not exist
	var dir = strings.Split(config.LogFilePath, "/")[0]
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	// create .log file if does not exist
	file, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	// create logger for environment
	var logger *zap.Logger
	if os.Getenv("APP_ENV") == "development" {
		// create a development encoder configuration
		encoderConfig := zap.NewDevelopmentEncoderConfig()

		encoderConfig.TimeKey = "timestamp"
		encoderConfig.LevelKey = "logLevel"
		encoderConfig.MessageKey = "message"
		encoderConfig.CallerKey = "caller"
		encoderConfig.StacktraceKey = "stack"
		encoderConfig.FunctionKey = zapcore.OmitKey
		encoderConfig.LineEnding = zapcore.DefaultLineEnding
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

		// create encoder for console
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zap.DebugLevel)

		// create encoder for file
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // For color output
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)

		// Combine the cores into a single logger
		logger = zap.New(zapcore.NewTee(consoleCore, fileCore))
		logger = logger.WithOptions(zap.AddCaller())

	} else {
		panic("todo")
	}

	return logger
}

// Section: main
func main() {
	// load config
	config := loadConfig()

	// setup logger as set globally
	logger := setupLogger(config)
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	zap.L().Debug("configuration values set", zap.String("config", fmt.Sprintf("%+v", config)))

	// init sarama config
	saramaConfig := sarama.NewConfig()
	client := createKafkaClient(config.BrokerAddresses, saramaConfig)

	// if topic does not exist, create
	exists := topicExists(client, config.TopicName)
	if !exists {
		createTopic(config, client, config.TopicName)
	} else {
		zap.L().Info("Topic already exists.", zap.String("topicName", config.TopicName))
	}
}
