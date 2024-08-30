package lib

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupZapLogger(config Config) (*zap.Logger, error) {
	// create dir if does not exist
	var dir = strings.Split(config.LogFilePath, "/")[0]
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	// create .log file if does not exist
	file, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error opening log file: '%v'", err)
	}

	// create logger for environment
	var logger *zap.Logger
	if config.AppEnv == "development" {
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
		logger = zap.New(zapcore.NewTee(consoleCore, fileCore)).WithOptions(zap.AddCaller())

	} else {
		fmt.Print("Invalid environment variable.")
	}

	return logger, nil
}
