package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

func New() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("can't create logger instance")
		os.Exit(1)
	}
	return logger
}
