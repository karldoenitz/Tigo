// tiger插件，一个脚手架工具，用于来初始化一个Tigo项目
package main

import "os"

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
	cmdVerbose = `
use command tiger to create a Tigo projection.

Usage:

    tiger <command> [args]

The commands are:

    create          to create a Tigo projection
    conf            to add a configuration for Tigo projection
    addHandler      to add a handler for Tigo projection

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

func main() {
	// TODO 此处补充脚手架逻辑
	// 获取命令行参数，根据参数判断是否是创建demo，
	// 如果创建demo，则直接把`DemoCode`注入到目标文件中就行
	// tiger支持的命令:
	//  - create xxx: 创建项目
	//  - addHandler xxx: 增加xxx命名的handler
	//  - conf xxx: 用xxx命名的配置文件替换现有配置文件，没有则新建

}
