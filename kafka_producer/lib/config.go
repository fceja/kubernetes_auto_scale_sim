package lib

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/joho/godotenv"
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

// Checks if reflect values are empty
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		// check if the slice length is 0
		if v.Len() == 0 {
			return true
		}

		// for elements in slice, check if contains zero values
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if !elem.IsZero() {
				return false
			}
		}
		return true
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}

// Validates config values are not empty
// Helps maintain consistency for local and docker config
func validateConfig(config Config) {
	configReflect := reflect.ValueOf(config)
	configReflectTyp := reflect.TypeOf(config)

	for i := 0; i < configReflect.NumField(); i++ {
		nilConfigFieldType := configReflectTyp.Field(i).Name
		configField := configReflect.Field(i)

		// check if config fields are missing or empty
		if isEmptyValue(configField) {
			panic(fmt.Sprintf("\nMissing or empty config env var: '%v'", nilConfigFieldType))
		}
	}
}

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
		LogLevel:        os.Getenv("LOG_LEVEL"),
		LogFilePath:     os.Getenv("LOG_FILE_PATH"),
		MessageLimit:    ConvertStrToInt(os.Getenv("MESSAGE_LIMIT")),
		SleepTimeout:    time.Duration(ConvertStrToInt(os.Getenv("SLEEP_TIMEOUT"))) * time.Millisecond,
		TopicName:       os.Getenv("TOPIC_NAME"),
	}
}

// Looks for Docker specific file to determine
// if running in a docker container
func isRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

// Looks for Kubernetes specific env var to determine
// if running in a Kubernetes cluster
func isRunningInKubernetes() bool {
	return os.Getenv("KUBERNETES_SERVICE_HOST") != ""
}

// Section: Main
func LoadConfig() Config {
	var config Config

	isRunningInDocker := isRunningInDocker()
	isRunningInKubernetes := isRunningInKubernetes()

	if !isRunningInDocker && !isRunningInKubernetes {
		// load .env file
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
		config = createConfig(false)

		} else if isRunningInKubernetes {
		fmt.Print("Is running in Kubernetes cluster.\n")
		config = createConfig(false)

	} else if isRunningInDocker {
		fmt.Print("Is running in docker.\n")
		config = createConfig(true)

	} else {
		panic("Logic error.")
	}

	// validate config
	validateConfig(config)

	return config
}
