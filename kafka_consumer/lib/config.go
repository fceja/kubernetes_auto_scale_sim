package lib

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	BrokerAddresses []string
	ConsumerGroupId string
	LogLevel        string
	LogFilePath     string
	TopicName       string
}

// Helper

// Section: Load env vars for Docker or Local
func createConfig(isDocker bool) Config {
	var brokerAddresses []string

	if isDocker {
		brokerAddresses = strings.Split(os.Getenv("DOCKER_BROKER_ADDRESSES"), ",")
	} else {
		brokerAddresses = strings.Split(os.Getenv("LOCAL_BROKER_ADDRESSES"), ",")
	}
	return Config{
		AppEnv:          os.Getenv("APP_ENV"),
		BrokerAddresses: brokerAddresses,
		ConsumerGroupId: os.Getenv("CONSUMER_GROUP_ID"),
		LogLevel:        os.Getenv("LOG_LEVEL"),
		LogFilePath:     os.Getenv("LOG_FILE_PATH"),
		TopicName:       os.Getenv("TOPIC_NAME"),
	}
}

// Validates config values are not empty
// Helps maintain consistency for local and docker config
func validateConfig(config Config) {
	value := reflect.ValueOf(config)
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).IsZero() {
			typ := reflect.TypeOf(config)
			panic(fmt.Sprintln("Invalid value:", typ.Field(i).Name))
		}
	}
}

// Looks for docker specific file to determine
// if running in a docker container
func isRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		fmt.Print("\nis running in docker.\n")
		// exists, running in docker
		return true
	}
	// does not exist, running locally
	fmt.Print("\nis NOT running in docker.\n")
	return false
}

// Section: Main
func LoadConfig() Config {
	var config Config

	if !isRunningInDocker() {
		// load .env file
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
		config = createConfig(false)
	} else {
		config = createConfig(true)
	}

	// validate config
	validateConfig(config)

	return config
}
