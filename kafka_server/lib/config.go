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
	LogFilePath     string
}

// Section: Helpers
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

// Create env vars config
func createConfig() Config {
	return Config{
		AppEnv:          os.Getenv("APP_ENV"),
		BrokerAddresses: strings.Split(os.Getenv("BROKER_ADDRESSES"), ","),
		LogFilePath:     os.Getenv("LOG_FILE_PATH"),
	}
}

// Validate config
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
	}

	// create and validate config
	config = createConfig()
	validateConfig(config)

	return config
}
