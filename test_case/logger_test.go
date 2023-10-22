package test_case

import (
	"Tigo/logger"
	"testing"
)

func TestInitLoggerWithObject(t *testing.T) {
	logLevel := logger.LogLevel{}
	logger.InitLoggerWithObject(logLevel)
	logger.Info.Println("log test passed")
}
