package test_case

import (
	"Tigo/logger"
	"testing"
)

func TestInitLoggerWithObject(t *testing.T) {
	logLevel := logger.LogLevel{}
	// 这里要增加相关属性，在终端命令行输出日志
	logger.InitLoggerWithObject(logLevel)
	logger.Info.Println("log test passed")
}
