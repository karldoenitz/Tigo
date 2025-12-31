package main

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
	"{{ .ProjectName }}/handler"
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

type {{ .HandlerName }} struct {
	web.BaseHandler
}

func (p *{{ .HandlerName }}) Get() {
	// write your code here
	p.ResponseAsText("Pong")
}

func (p *{{ .HandlerName }}) Post() {
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
	"cookie": "{{ .CookieKey }}",
	"ip": "0.0.0.0",
	"port": 8080,
	"log": {
		"trace": "stdout",
		"info": "{{ .WorkDir }}/log/tigo-framework-info.log",
		"warning": "{{ .WorkDir }}/log/tigo-framework-warning.log",
		"error": "{{ .WorkDir }}/log/tigo-framework-info-error.log"
	}
}
`
	configCodeYaml = `cookie: {{ .CookieKey }}
ip: 0.0.0.0
port: 8080
log:
  trace: stdout
  info: "{{ .WorkDir }}/log/tigo-framework-info.log"
  warning: "{{ .WorkDir }}/log/tigo-framework-warning.log"
  error: "{{ .WorkDir }}/log/tigo-framework-info-error.log"
`
	cmdVerbose = `
use command tiger to create a Tigo projection.

Usage:

    tiger <command> [args]

The commands are:

    addhandler      to add a handler for Tigo projection
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
