// Copyright 2018 The Tigo Authors. All rights reserved.
package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"encoding/json"
)

var (
	Trace        *log.Logger
	Info         *log.Logger
	Warning      *log.Logger
	Error        *log.Logger
)

var logPath = ""

// log分级结构体
//   - Trace    跟踪
//   - Info     信息
//   - Warning  预警
//   - Error    错误
type LogLevel struct {
	Trace    string   `json:"trace"`
	Info     string   `json:"info"`
	Warning  string   `json:"warning"`
	Error    string   `json:"error"`
}

// log文件路径与文件对象的关系映射
var logFileMapping = map[string] *os.File{}

// 更新log文件路径与log文件对象的映射关系
func updateLogMapping(filePath string) {
	if filePath != "" && filePath != "discard" && filePath != "stdout" {
		_, isExist := logFileMapping[filePath]
		if !isExist {
			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalln("Failed to open error log file: ", err)
				panic("Open File Error!")
			}
			logFileMapping[filePath] = file
		}
	}
}

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

// 设置log输出路径，警告：若使用了InitLoggerWithConfigFile和InitLoggerWithObject请不要使用此方法，会覆盖原有的log输出结构。
func SetLogPath(defineLogPath string) {
	logPath = defineLogPath
	initLogger()
}

// 根据配置文件路径初始化log模块；
// 配置文件需要配置如下部分：
//   - trace    "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
//   - info     "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
//   - warning  "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
//   - error    "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
func InitLoggerWithConfigFile(filePath string) {
	if filePath == "" {
		return
	}
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		os.Exit(1)
	}
	logLevel := LogLevel{}
	json.Unmarshal(raw, &logLevel)
	InitLoggerWithObject(logLevel)
}

// 根据LogLevel结构体的实例初始化log模块；
// 配置文件需要配置如下部分：
//   - Trace    "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
//   - Info     "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
//   - Warning  "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
//   - Error    "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
func InitLoggerWithObject(logLevel LogLevel)  {
	updateLogMapping(logLevel.Trace)
	updateLogMapping(logLevel.Info)
	updateLogMapping(logLevel.Warning)
	updateLogMapping(logLevel.Error)
	InitTrace(logLevel.Trace)
	InitInfo(logLevel.Info)
	InitWarning(logLevel.Warning)
	InitError(logLevel.Error)
}

// 初始化Trace，默认情况下不输出
func InitTrace(level string) {
	switch {
	case level == "" || level == "discard":
		Trace = log.New(ioutil.Discard, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	case level == "stdout":
		Trace = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	default:
		logFile := logFileMapping[level]
		Trace = log.New(io.MultiWriter(logFile, os.Stderr), "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// 初始化Info，默认情况下输出到终端
func InitInfo(level string)  {
	switch {
	case level == "" || level == "discard":
		Info = log.New(ioutil.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	case level == "stdout":
		Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	default:
		logFile := logFileMapping[level]
		Info = log.New(io.MultiWriter(logFile, os.Stderr), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	}
}

// 初始化Warning，默认情况下输出到终端
func InitWarning(level string)  {
	switch {
	case level == "" || level == "discard":
		Warning = log.New(ioutil.Discard, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	case level == "stdout":
		Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	default:
		logFile := logFileMapping[level]
		Warning = log.New(io.MultiWriter(logFile, os.Stderr), "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	}
}

// 初始化Warning，默认情况下输出到文件
func InitError(level string)  {
	switch {
	case level == "" || level == "discard":
		Error = log.New(ioutil.Discard, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	case level == "stdout":
		Error = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	default:
		logFile := logFileMapping[level]
		Error = log.New(io.MultiWriter(logFile, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	}
}
