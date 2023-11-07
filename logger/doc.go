// Package logger
// logger 包是Tigo框架的日志组件包，内部里面封装了日志相关操作，包括从配置文件初始化日志，通过config对象初始化日志，以及配置日志等级等。
// 目前该包不支持分布式日志，如果要使用分布式日志，推荐使用其他的第三方日志模块。
// ---------------------------------------------------------------------------------------------------------------------
// 初始化日志配置，如下示例所示：
//
// Basic Example:
//
//	logLevel := logger.LogLevel{
//		Info: "stdout",
//		Warning: "/path/to/logfile",
//		Error: "/path/to/logfile",
//		Trace: "/path/to/logfile",
//		TimeRoll: "D*3",
//	}
//	logger.InitLoggerWithObject(logLevel)
//
// ---------------------------------------------------------------------------------------------------------------------
// 日志打印用例：
//
// Basic Example:
//
//	logger.Info.Print("Hello Logger")
//	logger.Info.Println("Hello Logger")
//	logger.Info.Printf("Hello %d", 250)
package logger
