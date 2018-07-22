// Copyright 2018 The Tigo Authors. All rights reserved.
package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace        *log.Logger
	Info         *log.Logger
	Warning      *log.Logger
	Error        *log.Logger
)

var logPath = ""

// 初始化log模块
func initLogger() {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file: ", err)
	}
	Trace = log.New(ioutil.Discard, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	// 将运行日志写入控制台
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	// 将错误日志写入log文件
	Error = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// 初始化函数，加载log模块时运行
func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	logPath = dir + "/log.log"
	initLogger()
}

// 设置log输出路径
func SetLogPath(defineLogPath string) {
	logPath = defineLogPath
	initLogger()
}
