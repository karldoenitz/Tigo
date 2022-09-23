// tiger插件，一个脚手架工具，用于来初始化一个Tigo项目
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	DemoCode = `
package main

import (
	"github.com/karldoenitz/Tigo/TigoWeb"
)

// HelloHandler it's a demo handler
type HelloHandler struct {
    TigoWeb.BaseHandler
}

// Get http get method
func (h *HelloHandler) Get() {
	// write your code here
	h.ResponseAsHtml("<p1 style='color: red'>Hello Tiger Go!</p1>")
}

// urls url mapping
var urls = []TigoWeb.Pattern{
	{"/hello-world", HelloHandler{}, nil},
}

func main() {
	application := TigoWeb.Application{
		IPAddress:   "0.0.0.0",
		Port:        8888,
		UrlPatterns: urls,
	}
	application.Run()
}
`
	mainCode = `
package main

import (
	"github.com/karldoenitz/Tigo/TigoWeb"
	"%s/handler"
)

// Write you url mapping here
var urls = []TigoWeb.Pattern{
	{"/ping", handler.PingHandler{}, nil},
}

func main() {
	application := TigoWeb.Application{
		IPAddress:   "0.0.0.0",
		Port:        8080,
		UrlPatterns: urls,
	}
	application.Run()
}

`
	handlerCode = `
package handler

import (
	"github.com/karldoenitz/Tigo/TigoWeb"
)

type PingHandler struct {
	TigoWeb.BaseHandler
}

func (p *PingHandler) Get() {
	p.ResponseAsText("Pong")
}

// you can write your code here.
// to add 'Post', 'Put', 'Delete' and other methods here.

`
	cmdVerbose = `
use command tiger to create a Tigo projection.

Usage:

    tiger <command> [args]

The commands are:

    create          to create a Tigo projection
    conf            to add a configuration for Tigo projection
    addHandler      to add a handler for Tigo projection
    mod             to run go mod

Use "go help <command>" for more information about a command.

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
//  - args: 命令行输入的参数
func printCmdUsage(args []string) {
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
//  - args: 执行参数
func execEngine(args []string) {
	switch args[0] {
	case "create":
		execCreate(args[1])
		break
	case "conf":
		break
	case "addHandler":
		break
	}
}

// execCreate 执行create命令
//  - arg create命令的参数
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
	}
	_, _ = fHandler.WriteString(handlerCode)
	fmt.Printf("project `%s` created successfully\n", arg)
	fmt.Println("Execute go mod")
}

// execCmd 执行cmd命令
//  - commands: 需要执行的命令
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
	execCmd([]string{"go", "mod", "init", proName})
	execCmd([]string{"go", "mod", "tidy"})
	execCmd([]string{"go", "mod", "vendor"})
}

func main() {
	// 获取命令行参数，根据参数判断是否是创建demo，
	// 如果创建demo，则直接把`DemoCode`注入到目标文件中就行
	// tiger支持的命令:
	//  - create xxx: 创建项目
	//  - addHandler xxx: 增加xxx命名的handler
	//  - conf xxx: 用xxx命名的配置文件替换现有配置文件，没有则新建
	//  - mod: 进行go mod
	args, argsCnt := getCmdArgs()
	if argsCnt < 1 {
		fmt.Print(cmdVerbose)
		return
	}
	if args[0] == "mod" {
		goMod()
		return
	}
	if args[0] == "help" {
		printCmdUsage(args)
		return
	}
	execEngine(args)
}
