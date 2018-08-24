// Copyright 2018 The Tigo Authors. All rights reserved.
package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"encoding/json"
	"strings"
	"gopkg.in/yaml.v2"
	"fmt"
)

////////////////////////////////////////////////////常量/////////////////////////////////////////////////////////////////

var (
	Trace        *TiLog
	Info         *TiLog
	Warning      *TiLog
	Error        *TiLog
)

const (
	TraceLevel int = iota + 1
    InfoLevel
	WarningLevel
	ErrorLevel
)

var logPath = ""

var formater = map[int] string {
	TraceLevel:   "\x1b[32m %s \x1b[0m",
	InfoLevel:    "\x1b[34m %s \x1b[0m",
	WarningLevel: "\x1b[33m %s \x1b[0m",
	ErrorLevel:   "\x1b[31m %s \x1b[0m",
}

////////////////////////////////////////////////////结构体///////////////////////////////////////////////////////////////

// log分级结构体
//   - Trace    跟踪
//   - Info     信息
//   - Warning  预警
//   - Error    错误
// discard: 丢弃，stdout: 终端输出，文件路径表示log具体输出的位置
type LogLevel struct {
	Trace    string   `json:"trace"`
	Info     string   `json:"info"`
	Warning  string   `json:"warning"`
	Error    string   `json:"error"`
}

type TiLog struct {
	*log.Logger
	Level int
}

func (l *TiLog) Printf(format string, v ...interface{}) {
	formatStr := formater[l.Level]
	format = fmt.Sprintf(formatStr, format)
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *TiLog) Print(v ...interface{}) {
	formatStr := formater[l.Level]
	logInfo := fmt.Sprintf(formatStr, fmt.Sprint(v...))
	l.Output(2, logInfo)
}

func (l *TiLog) Println(v ...interface{}) {
	formatStr := formater[l.Level]
	logInfo := fmt.Sprintf(formatStr, fmt.Sprintln(v...))
	l.Output(2, logInfo)
}

////////////////////////////////////////////////////初始化logger的方法集//////////////////////////////////////////////////

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
	Trace = &TiLog{}
	Trace.Logger = log.New(ioutil.Discard, "\x1b[32m TRACE:   \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	Trace.Level = TraceLevel
	// 将运行日志写入控制台
	Info = &TiLog{}
	Info.Logger = log.New(os.Stdout, "\x1b[34m INFO:    \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	Info.Level = InfoLevel
	Warning = &TiLog{}
	Warning.Logger = log.New(os.Stdout, "\x1b[33m WARNING: \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning.Level = WarningLevel
	// 将错误日志写入log文件
	Error = &TiLog{}
	Error.Logger = log.New(io.MultiWriter(file, os.Stderr), "\x1b[31m ERROR:   \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	Error.Level = ErrorLevel
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
	if strings.HasSuffix(filePath, ".json") {
		json.Unmarshal(raw, &logLevel)
	}
	if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, "yml") {
		yaml.Unmarshal(raw, &logLevel)
	}
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
	Trace.Level = TraceLevel
	switch {
	case level == "" || level == "discard":
		Trace.Logger = log.New(ioutil.Discard, "\x1b[32m TRACE:   \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	case level == "stdout":
		Trace.Logger = log.New(os.Stdout, "\x1b[32m TRACE:   \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	default:
		logFile := logFileMapping[level]
		Trace.Logger = log.New(io.MultiWriter(logFile, os.Stderr), "\x1b[32m TRACE:   \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

// 初始化Info，默认情况下输出到终端
func InitInfo(level string)  {
	Info.Level = InfoLevel
	switch {
	case level == "" || level == "discard":
		Info.Logger = log.New(ioutil.Discard, "\x1b[34m INFO:    \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	case level == "stdout":
		Info.Logger = log.New(os.Stdout, "\x1b[34m INFO:    \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	default:
		logFile := logFileMapping[level]
		Info.Logger = log.New(io.MultiWriter(logFile, os.Stderr), "\x1b[34m INFO:    \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	}
}

// 初始化Warning，默认情况下输出到终端
func InitWarning(level string)  {
	Warning.Level = WarningLevel
	switch {
	case level == "" || level == "discard":
		Warning.Logger = log.New(ioutil.Discard, "\x1b[33m WARNING: \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	case level == "stdout":
		Warning.Logger = log.New(os.Stdout, "\x1b[33m WARNING: \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	default:
		logFile := logFileMapping[level]
		Warning.Logger = log.New(io.MultiWriter(logFile, os.Stderr), "\x1b[33m WARNING: \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	}
}

// 初始化Warning，默认情况下输出到文件
func InitError(level string)  {
	Error.Level = ErrorLevel
	switch {
	case level == "" || level == "discard":
		Error.Logger = log.New(ioutil.Discard, "\x1b[31m ERROR:   \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	case level == "stdout":
		Error.Logger = log.New(os.Stdout, "\x1b[31m ERROR:   \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	default:
		logFile := logFileMapping[level]
		Error.Logger = log.New(io.MultiWriter(logFile, os.Stderr), "\x1b[31m ERROR:   \x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
		break
	}
}
