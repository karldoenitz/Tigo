// tiger插件，一个脚手架工具，用于来初始化一个Tigo项目
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/karldoenitz/Tigo/web"
)

const (
	DemoCode = `package main

import (
	"github.com/karldoenitz/Tigo/web"
)

// HelloHandler it's a demo handler
type HelloHandler struct {
    web.BaseHandler
}

// Get http get method
func (h *HelloHandler) Get() {
	// write your code here
	h.ResponseAsHtml("<p1 style='color: red'>Hello Tiger Go!</p1>")
}

// urls url mapping
var urls = []web.Pattern{
	{"/hello-world", HelloHandler{}, nil},
}

func main() {
	application := web.Application{
		IPAddress:   "0.0.0.0",
		Port:        8888,
		UrlPatterns: urls,
	}
	application.Run()
}
`
	mainCode = `package main

import (
	"github.com/karldoenitz/Tigo/web"
	"%s/handler"
)

// Write you url mapping here
var urls = []web.Pattern{
	{"/ping", handler.PingHandler{}, nil},
}

func main() {
	application := web.Application{
		IPAddress:   "0.0.0.0",
		Port:        8080,
		UrlPatterns: urls,
	}
	application.Run()
}

`
	handlerCode = `// you can write your code here.
// You can add 'Post', 'Put', 'Delete' and other methods to handler.
package handler

import (
	"github.com/karldoenitz/Tigo/web"
)

type %s struct {
	web.BaseHandler
}

func (p *%s) Get() {
	// write your code here
	p.ResponseAsText("Pong")
}

func (p *%s) Post() {
	// write your code here
	p.ResponseAsText("Pong")
}

`
	logCode = `// you can write your code here.
// You can modify the log level and add more logs.
package logger

import (
	"os"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)
}
`
	configCodeJson = `{
	"cookie": "%s",
	"ip": "0.0.0.0",
	"port": 8080,
	"log": {
		"trace": "stdout",
		"info": "%s/log/tigo-framework-info.log",
		"warning": "%s/log/tigo-framework-warning.log",
		"error": "%s/log/tigo-framework-info-error.log"
	}
}
`
	configCodeYaml = `cookie: %s
ip: 0.0.0.0
port: 8080
log:
  trace: stdout
  info: "%s/log/tigo-framework-info.log"
  warning: "%s/log/tigo-framework-warning.log"
  error: "%s/log/tigo-framework-info-error.log"
`
	cmdVerbose = `
use command tiger to create a Tigo projection.

Usage:

    tiger <command> [args]

The commands are:

    addHandler      to add a handler for Tigo projection
    create          to create a Tigo projection
    conf            to add a configuration for Tigo projection
    logger          to add a logger for Tigo projection
    mod             to run go mod
    version         to show Tigo version

Use "tiger help <command>" for more information about a command.

`
	cmdCreateVerbose = `
use this command to create a Tigo project.
"tiger create <project_name>" can create a project with name "project_name",
"tiger create demo" can create a demo project.

`
	cmdConfVerbose = `
use this command to add a configuration.
if it's an empty folder, this command will throw an error.
the new configuration will replace the old configuration.

`
	cmdAddHandlerVerbose = `
use this command to add a handler with defined name.
"tiger addHandler <handler_name>" will add a handler named "handler_name".

`
)

// getWorkingDirPath 获取当前工作路径
func getWorkingDirPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// getCmdArgs 获取命令行参数及命令个数
func getCmdArgs() (args []string, argNum int) {
	args = os.Args[1:]
	argNum = len(args)
	return
}

// printCmdUsage 打印help
//   - args: 命令行输入的参数
func printCmdUsage(args []string) {
	args = append(args, "")
	cmd := args[1]
	switch cmd {
	case "create":
		fmt.Print(cmdCreateVerbose)
		break
	case "conf":
		fmt.Print(cmdConfVerbose)
		break
	case "addHandler":
		fmt.Print(cmdAddHandlerVerbose)
		break
	default:
		fmt.Print(cmdVerbose)
		break
	}
}

// execEngine 执行引擎
//   - args: 执行参数
func execEngine(args []string) {
	switch args[0] {
	case "create":
		execCreate(args[1])
		break
	case "conf":
		execConf(args[1])
		break
	case "addHandler":
		execAddHandler(args[1])
		break
	case "logger":
		execAddLogger()
		break
	}
}

// execCreate 执行create命令
//   - arg create命令的参数
func execCreate(arg string) {
	// 先创建目录
	workDir := getWorkingDirPath()
	projectPath := fmt.Sprintf("%s/%s", workDir, arg)
	if err := os.Mkdir(projectPath, os.ModePerm); err != nil {
		panic(err.Error())
	}
	if arg == "demo" {
		// 再创建文件
		f, err := os.Create(fmt.Sprintf("%s/main.go", projectPath))
		if err != nil {
			panic(err.Error())
		}
		if _, err := f.WriteString(DemoCode); err != nil {
			panic(err)
		}
		fmt.Println("project `demo` created successfully")
		fmt.Println("Execute go mod")
		_ = f.Close()
		return
	}
	// 创建非demo项目的main文件
	f, err := os.Create(fmt.Sprintf("%s/main.go", projectPath))
	if err != nil {
		panic(err.Error())
	}
	if _, err := f.WriteString(fmt.Sprintf(mainCode, arg)); err != nil {
		panic(err)
	}
	// 创建handler文件
	if err := os.Mkdir(projectPath+"/handler", os.ModePerm); err != nil {
		fmt.Println(err.Error())
	}
	fHandler, err := os.Create(fmt.Sprintf("%s/handler/pinghandler.go", projectPath))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, _ = fHandler.WriteString(fmt.Sprintf(handlerCode, "PingHandler", "PingHandler", "PingHandler"))
	_ = f.Close()
	_ = fHandler.Close()

	fmt.Printf("project `%s` created successfully\n", arg)
	fmt.Println("Execute go mod")
}

// execCmd 执行cmd命令
//   - commands: 需要执行的命令
func execCmd(commands []string) bool {
	cmd := exec.Command(commands[0], commands[1:]...)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
		return false
	}
	return true
}

// goMod 执行go mod
func goMod() {
	dir, _ := os.Getwd()
	splitPath := strings.Split(dir, "/")
	proName := splitPath[len(splitPath)-1]
	_ = os.Setenv("GO111MODULE", "on")
	execCmd([]string{"go", "mod", "init", proName})
	execCmd([]string{"go", "mod", "tidy"})
	execCmd([]string{"go", "mod", "vendor"})
}

// execAddHandler 在当前Tigo项目中增加一个handler
//   - handlerName: handler名字
func execAddHandler(handlerName string) {
	workDir := getWorkingDirPath()
	handlerPath := fmt.Sprintf("%s/handler", workDir)
	_ = os.Mkdir(handlerPath, os.ModePerm)
	// 如果有则新建一个handler文件，并注入代码
	fileName := strings.ToLower(handlerName)
	fHandler, err := os.Create(fmt.Sprintf("%s/%s.go", handlerPath, fileName))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, _ = fHandler.WriteString(fmt.Sprintf(handlerCode, handlerName, handlerName, handlerName))
	_ = fHandler.Close()
	// 判断是否有 main 文件
	_, err = os.Stat(fmt.Sprintf("%s/main.go", workDir))
	if err != nil {
		// 如果没有则退出
		fmt.Println(err.Error())
		return
	}
	// 如果有则检测代码，并在urls中插入一个url映射
	content, err := os.ReadFile(fmt.Sprintf("%s/main.go", workDir))
	if err != nil {
		fmt.Printf("read file error:%v\n", err)
		return
	}
	// 寻找main.go中的url配置
	codes := strings.Split(string(content), "\n")
	var isFoundUrls bool
	var newCodes []string
	url := strings.Replace(fileName, "handler", "", -1)
	for _, code := range codes {
		if code == "var urls = []web.Pattern{" {
			isFoundUrls = true
		}
		if code == "}" && isFoundUrls {
			code = fmt.Sprintf("\t{\"/%s\", handler.%s{}, nil},\n}", url, handlerName)
			isFoundUrls = false
			newCodes = append(newCodes, code)
			continue
		}
		newCodes = append(newCodes, code)
	}
	newCode := strings.Join(newCodes, "\n")
	f, err := os.Create(fmt.Sprintf("%s/main.go", workDir))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, _ = f.WriteString(newCode)
	_ = f.Close()
}

// execAddLogger 增加日志配置
func execAddLogger() {
	// TODO 这里增加logrus配置，后续增加支持，目前只实现一个入口，具体逻辑还没添加，工作太忙了，没空维护
	workDir := getWorkingDirPath()
	loggerPath := fmt.Sprintf("%s/common", workDir)
	_ = os.Mkdir(loggerPath, os.ModePerm)
	// 如果有则新建一个handler文件，并注入代码
	lHandler, err := os.Create(fmt.Sprintf("%s/logger.go", loggerPath))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, _ = lHandler.WriteString(fmt.Sprintf(logCode))
	_ = lHandler.Close()
}

// execConf 增加配置文件
//   - arg: 配置文件名称
func execConf(arg string) {
	workDir := getWorkingDirPath()
	configPath := fmt.Sprintf("%s/%s", workDir, arg)
	_ = os.Mkdir(fmt.Sprintf("%s/log", workDir), os.ModePerm)
	f, err := os.Create(configPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	currentTime := time.Now().String() + arg
	if strings.HasSuffix(arg, ".json") {
		_, _ = f.WriteString(fmt.Sprintf(configCodeJson, web.MD5m16(currentTime), workDir, workDir, workDir))
	} else {
		_, _ = f.WriteString(fmt.Sprintf(configCodeYaml, web.MD5m16(currentTime), workDir, workDir, workDir))
	}
	_ = f.Close()
	content, err := os.ReadFile(fmt.Sprintf("%s/main.go", workDir))
	if err != nil {
		fmt.Printf("read file error:%v\n", err)
		return
	}
	// 寻找main.go中的application.Run()配置
	codes := strings.Split(string(content), "\n")
	var newCodes []string
	for _, code := range codes {
		if code == "\tapplication.Run()" {
			code = fmt.Sprintf("\tapplication.ConfigPath = \"%s\"\n\tapplication.Run()", configPath)
			newCodes = append(newCodes, code)
			continue
		}
		newCodes = append(newCodes, code)
	}
	newCode := strings.Join(newCodes, "\n")
	f, err = os.Create(fmt.Sprintf("%s/main.go", workDir))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = f.WriteString(newCode)
	if err != nil {
		fmt.Println(err.Error())
	}
	_ = f.Close()
}

func main() {
	// 获取命令行参数，根据参数判断是否是创建demo，
	// 如果创建demo，则直接把常变量`DemoCode`注入到目标文件中就行
	// tiger支持的命令:
	//  - create xxx: 创建项目
	//  - addHandler xxx: 增加xxx命名的handler
	//  - conf xxx: 用xxx命名的配置文件替换现有配置文件，没有则新建
	//  - mod: 进行go mod
	//  - version: 获取当前Tigo版本号
	args, argsCnt := getCmdArgs()
	if argsCnt < 1 {
		fmt.Print(cmdVerbose)
		return
	}
	if args[0] == "mod" {
		goMod()
		return
	}
	if args[0] == "version" {
		fmt.Println(web.Version)
		return
	}
	if args[0] == "help" {
		printCmdUsage(args)
		return
	}
	execEngine(args)
}
