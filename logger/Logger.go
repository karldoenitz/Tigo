// Copyright 2018 The Tigo Authors. All rights reserved.
package logger

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

var dateFormatter = ".2006-01-02_15:04:05"

var logPath = ""

var formatter = map[int] string {
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
	TimeRoll string   `json:"time_roll"`
}

type TiLog struct {
	*log.Logger
	Level int
}

func (l *TiLog) Printf(format string, v ...interface{}) {
	formatStr := formatter[l.Level]
	format = fmt.Sprintf(formatStr, format)
	l.Output(2, fmt.Sprintf(format, v...))
}

func (l *TiLog) Print(v ...interface{}) {
	formatStr := formatter[l.Level]
	logInfo := fmt.Sprintf(formatStr, fmt.Sprint(v...))
	l.Output(2, logInfo)
}

func (l *TiLog) Println(v ...interface{}) {
	formatStr := formatter[l.Level]
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
		println(err.Error())
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
	go startTimer(sliceLog, logLevel)
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

// 开启定时器，传入日志切分函数，日至等级对象
func startTimer(F func(LogLevel), logLevel LogLevel) {
	// 如果没有设置日志切分，则直接返回，不设置定时任务
	if logLevel.TimeRoll == "" {
		return
	}
	var ch chan int
	//定时任务
	timeRollingFrequency := getTimeRollingFrequency(logLevel)
	ticker := time.NewTicker(timeRollingFrequency)
	go func() {
		for range ticker.C {
			F(logLevel)
		}
		ch <- 1
	}()
	<-ch
}

// 解析time_roll配置文件，获取定时任务运行频率
// time_roll格式必须为:
//   - D*intN: 每N天切分一次日志
//   - H*intN: 每N小时切分一次日志
//   - M*intN: 每N分钟切分一次日志
//   - S*intN: 每N秒钟切分一次日志
func getTimeRollingFrequency(logLevel LogLevel) time.Duration {
	rollingMap := map[string] time.Duration {
		"D": time.Hour * 24,
		"H": time.Hour,
		"M": time.Minute,
		"S": time.Second,
	}
	dateMap := map[string] string {
		"D": ".2006-01-02",
		"H": ".2006-01-02_15",
		"M": ".2006-01-02_15:04",
		"S": ".2006-01-02_15:04:05",
	}
	rollingStr := logLevel.TimeRoll
	rollingArr := strings.Split(rollingStr, "*")
	if len(rollingArr) != 2 {
		println("time rolling string error")
		os.Exit(1)
	}
	timeFlag, isOk := rollingMap[rollingArr[0]]
	if !isOk {
		println("time rolling flag error")
		os.Exit(1)
	}
	timeOffset, err := strconv.Atoi(rollingArr[1])
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	// 设置格式化日期字符串
	dateFormatter = dateMap[rollingArr[0]]
	return timeFlag * time.Duration(timeOffset)
}

// 日志切分函数
func sliceLog(logLevel LogLevel)  {
	// 获取上一个切分节点的日志对象
	TraceLogFile, isTraceExisted := logFileMapping[logLevel.Trace]
	InfoLogFile, isInfoExisted := logFileMapping[logLevel.Info]
	WarningLogFile, isWarningExisted := logFileMapping[logLevel.Warning]
	ErrorLogFile, isErrorExisted := logFileMapping[logLevel.Error]
	// 获取当前切分节点的日志名称
	nowTimeStr := time.Now().Format(dateFormatter)
	logLevel.Trace += nowTimeStr
	logLevel.Info += nowTimeStr
	logLevel.Warning += nowTimeStr
	logLevel.Error += nowTimeStr
	updateLogMapping(logLevel.Trace)
	updateLogMapping(logLevel.Info)
	updateLogMapping(logLevel.Warning)
	updateLogMapping(logLevel.Error)
	// 如果上一份日志文件没有关闭则进行关闭
	if isTraceExisted && TraceLogFile != nil {
		TraceLogFile.Close()
	}
	if isInfoExisted && InfoLogFile != nil {
		InfoLogFile.Close()
	}
	if isWarningExisted && WarningLogFile != nil {
		WarningLogFile.Close()
	}
	if isErrorExisted && ErrorLogFile != nil {
		ErrorLogFile.Close()
	}
	// 重新初始化当前日志
	InitTrace(logLevel.Trace)
	InitInfo(logLevel.Info)
	InitWarning(logLevel.Warning)
	InitError(logLevel.Error)
}
