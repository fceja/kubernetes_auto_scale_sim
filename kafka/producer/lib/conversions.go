package lib

import (
	"strconv"

	"go.uber.org/zap"
)

func ConvertStrToInt(inputStr string) int {
	inputInt, err := strconv.Atoi(inputStr)
	if err != nil {
		zap.L().Fatal("Error converting string to int.", zap.Error(err))
	}

	return inputInt
}
