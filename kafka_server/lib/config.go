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
	BrokerAddresses []string // set by LOCAL_DOCKER_BROKER_ADDRESSES OR DOCKER_DOCKER_BROKER_ADDRESSES
	LogFilePath     string
}

// Section: Helpers
// Converts string to a map
func convertStrToMap(secretStr string) map[string]string {
	lines := strings.Split(secretStr, "\n")
	newMap := make(map[string]string)

	// extract key-value pairs
	for _, line := range lines {
		// skip empty lines
		if line == "" {
			continue
		}

		// split into key and value, by '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			panic(fmt.Sprintln("\nMalformed line:", line))
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		newMap[key] = value
	}

	return newMap
}

// Checks if reflect values are empty
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		// Check if the slice length is 0
		if v.Len() == 0 {
			return true
		}

		// If the slice length is 1, check if it contains an empty value
		if v.Len() == 1 {
			elem := v.Index(0) // Get the first element of the slice
			return elem.IsZero()
		}

		// For slices with more than one element, just check if all are zero values
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
// Load env vars via docker swarm - https://docs.docker.com/engine/swarm/
func loadDockerEnvConfig(envConfigPath string) Config {
	fmt.Println("Loading docker swarm config.")
	secretData, err := os.ReadFile(envConfigPath)
	if err != nil {
		panic(fmt.Sprintf("\nSwarm initialized? Stack enabled? Error: %v", err))
	}

	// convert to map
	mapSecrets := convertStrToMap(string(secretData))

	return Config{
		AppEnv:          mapSecrets["APP_ENV"],
		BrokerAddresses: strings.Split(mapSecrets["DOCKER_BROKER_ADDRESSES"], ","),
		LogFilePath:     mapSecrets["LOG_FILE_PATH"],
	}
}

// Load env vars from local .env file
func loadLocalEnvConfig() Config {
	fmt.Println("Loading local config.")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return Config{
		AppEnv:          os.Getenv("APP_ENV"),
		BrokerAddresses: strings.Split(os.Getenv("LOCAL_BROKER_ADDRESSES"), ","),
		LogFilePath:     os.Getenv("LOG_FILE_PATH"),
	}
}

// Section: Main
// First, determine whether running locally or within a docker container
// Then, load respective env vars
func LoadConfig() Config {
	var envConfig Config

	envConfigPath := ".env"
	_, envFileErr := os.ReadFile(envConfigPath)
	if envFileErr != nil {
		// load docker env vars from docker swarm
		const dockerSwarmSecretsPath = "/run/secrets/kafka-server-secrets"
		envConfig = loadDockerEnvConfig(dockerSwarmSecretsPath)
	} else {
		// load local env vars from .env file
		envConfig = loadLocalEnvConfig()
	}

	// validate config env vars are not empty
	validateConfig(envConfig)

	return envConfig
}
