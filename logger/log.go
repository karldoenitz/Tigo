// Package logger 提供Tigo框架自带的log纪录功能
package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// //////////////////////////////////////////////////常量/////////////////////////////////////////////////////////////////

// Trace 等变量不同级别的log实例
var (
	Trace   *TiLog
	Info    *TiLog
	Warning *TiLog
	Error   *TiLog
)

// TraceLevel 等变量表示log实例的级别
const (
	TraceLevel int = iota + 1
	InfoLevel
	WarningLevel
	ErrorLevel
)

var dateFormatter = ".2006-01-02_15:04:05"

var logPath = ""

var formatter = map[int]string{
	TraceLevel:   "\x1b[32m %s \x1b[0m",
	InfoLevel:    "\x1b[34m %s \x1b[0m",
	WarningLevel: "\x1b[33m %s \x1b[0m",
	ErrorLevel:   "\x1b[31m %s \x1b[0m",
}

// //////////////////////////////////////////////////结构体///////////////////////////////////////////////////////////////

// LogLevel 是log分级结构体
//   - Trace    跟踪
//   - Info     信息
//   - Warning  预警
//   - Error    错误
//
// discard: 丢弃，stdout: 终端输出，文件路径表示log具体输出的位置
type LogLevel struct {
	Trace    string `json:"trace" yaml:"trace"`
	Info     string `json:"info" yaml:"info"`
	Warning  string `json:"warning" yaml:"warning"`
	Error    string `json:"error" yaml:"error"`
	TimeRoll string `json:"time_roll" yaml:"timeRoll"`
}

// TiLog 是Tigo自定义的log结构体
type TiLog struct {
	*log.Logger
	Level         int
	consoleLogger *log.Logger
}

// Printf 格式化输出log
func (l *TiLog) Printf(format string, v ...interface{}) {
	formatStr := formatter[l.Level]
	format = fmt.Sprintf(formatStr, format)
	_ = l.Output(2, fmt.Sprintf(format, v...))
}

// Print 打印log，不换行
func (l *TiLog) Print(v ...interface{}) {
	formatStr := formatter[l.Level]
	logInfo := fmt.Sprintf(formatStr, fmt.Sprint(v...))
	_ = l.Output(2, logInfo)
}

// Println 打印log并且换行
func (l *TiLog) Println(v ...interface{}) {
	formatStr := formatter[l.Level]
	logInfo := fmt.Sprintf(formatStr, fmt.Sprintln(v...))
	_ = l.Output(2, logInfo)
}

// //////////////////////////////////////////////////初始化logger的方法集//////////////////////////////////////////////////

// log文件路径与文件对象的关系映射 值类型为*os.File
var logFileMapping = sync.Map{}

// 更新log文件路径与log文件对象的映射关系
func updateLogMapping(filePath string) {
	if filePath != "" && filePath != "discard" && filePath != "stdout" {
		_, isExist := logFileMapping.Load(filePath)
		if !isExist {
			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalln("Failed to open error log file: ", err)
			}
			logFileMapping.Store(filePath, file)
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
	Trace.Logger = log.New(io.MultiWriter(file), "TRACE   ", log.Ldate|log.Ltime)
	Trace.Level = TraceLevel
	Trace.consoleLogger = log.New(os.Stdout, "\x1b[42m TRACE   \x1b[0m ", log.Ldate|log.Ltime)
	Info = &TiLog{}
	Info.Logger = log.New(io.MultiWriter(file), "INFO    ", log.Ldate|log.Ltime)
	Info.Level = InfoLevel
	Info.consoleLogger = log.New(os.Stdout, "\x1b[44m INFO    \x1b[0m ", log.Ldate|log.Ltime)
	Warning = &TiLog{}
	Warning.Logger = log.New(io.MultiWriter(file), "WARNING ", log.Ldate|log.Ltime)
	Warning.Level = WarningLevel
	Warning.consoleLogger = log.New(os.Stdout, "\x1b[43m WARNING \x1b[0m ", log.Ldate|log.Ltime)
	// 将错误日志写入log文件
	Error = &TiLog{}
	Error.Logger = log.New(io.MultiWriter(file, os.Stderr), "ERROR   ", log.Ldate|log.Ltime)
	Error.Level = ErrorLevel
	Error.consoleLogger = log.New(os.Stdout, "\x1b[41m ERROR   \x1b[0m ", log.Ldate|log.Ltime)
}

// 初始化函数，加载log模块时运行
func init() {
	dir, err := os.Getwd()
	if err != nil {
		println(err.Error())
		dir = "/var"
	}
	logPath = dir + "/log.log"
	initLogger()
}

// InitLoggerWithObject 根据LogLevel结构体的实例初始化log模块；
// 配置文件需要配置如下部分：
//   - Trace    "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
//   - Info     "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
//   - Warning  "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
//   - Error    "discard": 不输出；"stdout": 终端输出不打印到文件；"/path/demo.log": 输出到指定文件
func InitLoggerWithObject(logLevel LogLevel) {
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

// InitTrace 初始化Trace，默认情况下不输出 TODO 后面的输出到控制台的log需要去掉，只保留写入到文件的日志
func InitTrace(level string) {
	Trace.Level = TraceLevel
	switch {
	case level == "" || level == "discard":
		Trace.Logger = log.New(io.Discard, "\x1b[42m TRACE   \x1b[0m ", log.Ldate|log.Ltime)
		break
	case level == "stdout":
		Trace.Logger = log.New(os.Stdout, "\x1b[42m TRACE   \x1b[0m ", log.Ldate|log.Ltime)
		break
	default:
		logFile, ok := logFileMapping.Load(level)
		if !ok {
			log.Print("Failed to open trace log file: ", level)
			break
		}
		Trace.Logger = log.New(io.MultiWriter(logFile.(*os.File), os.Stderr), "\x1b[42m TRACE   \x1b[0m ", log.Ldate|log.Ltime)
	}
}

// InitInfo 初始化Info，默认情况下输出到终端
func InitInfo(level string) {
	Info.Level = InfoLevel
	switch {
	case level == "" || level == "discard":
		Info.Logger = log.New(io.Discard, "\x1b[44m INFO    \x1b[0m ", log.Ldate|log.Ltime)
		break
	case level == "stdout":
		Info.Logger = log.New(os.Stdout, "\x1b[44m INFO    \x1b[0m ", log.Ldate|log.Ltime)
		break
	default:
		logFile, ok := logFileMapping.Load(level)
		if !ok {
			log.Print("Failed to open info log file: ", level)
			break
		}
		Info.Logger = log.New(io.MultiWriter(logFile.(*os.File), os.Stderr), "\x1b[44m INFO    \x1b[0m ", log.Ldate|log.Ltime)
		break
	}
}

// InitWarning 初始化Warning，默认情况下输出到终端
func InitWarning(level string) {
	Warning.Level = WarningLevel
	switch {
	case level == "" || level == "discard":
		Warning.Logger = log.New(io.Discard, "\x1b[43m WARNING \x1b[0m ", log.Ldate|log.Ltime)
		break
	case level == "stdout":
		Warning.Logger = log.New(os.Stdout, "\x1b[43m WARNING \x1b[0m ", log.Ldate|log.Ltime)
		break
	default:
		logFile, ok := logFileMapping.Load(level)
		if !ok {
			log.Print("Failed to open warning log file: ", level)
			break
		}
		Warning.Logger = log.New(io.MultiWriter(logFile.(*os.File), os.Stderr), "\x1b[43m WARNING \x1b[0m ", log.Ldate|log.Ltime)
		break
	}
}

// InitError 初始化Error，默认情况下输出到文件
func InitError(level string) {
	Error.Level = ErrorLevel
	switch {
	case level == "" || level == "discard":
		Error.Logger = log.New(io.Discard, "\x1b[41m ERROR   \x1b[0m ", log.Ldate|log.Ltime)
		break
	case level == "stdout":
		Error.Logger = log.New(os.Stdout, "\x1b[41m ERROR   \x1b[0m ", log.Ldate|log.Ltime)
		break
	default:
		logFile, ok := logFileMapping.Load(level)
		if !ok {
			log.Print("Failed to open error log file: ", level)
			break
		}
		Error.Logger = log.New(io.MultiWriter(logFile.(*os.File), os.Stderr), "\x1b[41m ERROR   \x1b[0m ", log.Ldate|log.Ltime)
		break
	}
}

// 开启定时器，传入日志切分函数，日至等级对象
func startTimer(F func(LogLevel, time.Time), logLevel LogLevel) {
	// 如果没有设置日志切分，则直接返回，不设置定时任务
	if logLevel.TimeRoll == "" {
		return
	}
	var ch chan int
	// 定时任务
	timeRollingFrequency := getTimeRollingFrequency(logLevel)
	ticker := time.NewTicker(timeRollingFrequency)
	go func() {
		for currentTime := range ticker.C {
			F(logLevel, currentTime)
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
	rollingMap := map[string]time.Duration{
		"D": time.Hour * 24,
		"H": time.Hour,
		"M": time.Minute,
		"S": time.Second,
	}
	dateMap := map[string]string{
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
func sliceLog(logLevel LogLevel, current time.Time) {
	// 获取上一个切分节点的日志对象
	TraceLogFile, isTraceExisted := logFileMapping.Load(logLevel.Trace)
	InfoLogFile, isInfoExisted := logFileMapping.Load(logLevel.Info)
	WarningLogFile, isWarningExisted := logFileMapping.Load(logLevel.Warning)
	ErrorLogFile, isErrorExisted := logFileMapping.Load(logLevel.Error)
	// 获取当前切分节点的日志名称
	nowTimeStr := current.Format(dateFormatter)
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
		_ = TraceLogFile.(*os.File).Close()
	}
	if isInfoExisted && InfoLogFile != nil {
		_ = InfoLogFile.(*os.File).Close()
	}
	if isWarningExisted && WarningLogFile != nil {
		_ = WarningLogFile.(*os.File).Close()
	}
	if isErrorExisted && ErrorLogFile != nil {
		_ = ErrorLogFile.(*os.File).Close()
	}
	// 重新初始化当前日志
	InitTrace(logLevel.Trace)
	InitInfo(logLevel.Info)
	InitWarning(logLevel.Warning)
	InitError(logLevel.Error)
}

// //////////////////////////////////////////////////http相关工具函数//////////////////////////////////////////////////////

// StatusColor 给http状态码进行终端着色
//   - status: 状态码
func StatusColor(status int) (coloredStatus string) {
	switch {
	case status < http.StatusMultipleChoices:
		coloredStatus = fmt.Sprintf("\x1B[6;30;32m[%d]\x1B[0m", status)
	case status < http.StatusBadRequest:
		coloredStatus = fmt.Sprintf("\x1B[6;30;34m[%d]\x1B[0m", status)
	case status < http.StatusInternalServerError:
		coloredStatus = fmt.Sprintf("\x1B[6;30;33m[%d]\x1B[0m", status)
	default:
		coloredStatus = fmt.Sprintf("\x1B[6;30;31m[%d]\x1B[0m", status)
	}
	return
}
