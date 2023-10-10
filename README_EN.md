[![Badge](https://img.shields.io/badge/link-Tigo-blue.svg)](https://karldoenitz.github.io/Tigo-EN/)
[![LICENSE](https://img.shields.io/badge/license-Tigo-blue.svg)](https://github.com/karldoenitz/Tigo/blob/master/LICENSE)
[![Go](https://github.com/karldoenitz/Tigo/actions/workflows/go.yml/badge.svg)](https://github.com/karldoenitz/Tigo/actions/workflows/go.yml)
[![Join the chat at https://gitter.im/karlooper/Tigo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/karlooper/Tigo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Open Source Helpers](https://www.codetriage.com/karldoenitz/tigo/badges/users.svg)](https://www.codetriage.com/karldoenitz/Tigo)
[![Go Report Card](https://goreportcard.com/badge/github.com/karldoenitz/Tigo)](https://goreportcard.com/report/github.com/karldoenitz/Tigo)
[![GoDoc](https://godoc.org/github.com/karldoenitz/Tigo?status.svg)](https://pkg.go.dev/github.com/karldoenitz/Tigo)
[![Release](https://img.shields.io/github/release/karldoenitz/Tigo.svg)](https://github.com/karldoenitz/Tigo/releases)  
![Tigo logo](https://raw.githubusercontent.com/karldoenitz/Tigo/master/documentation/tigo_logo.jpg "this is Tigo logo")
# Tigo([中文文档点击此处](https://github.com/karldoenitz/Tigo/blob/master/README.md))
A web framework developed in go language.

# Plugins and Tools for Tigo
- **tiger**  
`tiger` is a commandline tool for `Tigo` framework, you can use `tiger` to create a `Tigo` projection.  
[glance tiger](https://github.com/karldoenitz/tiger)  
- **tission**  
`tission` is a session plugin for `Tigo`.  
[glance tission](https://github.com/karldoenitz/tission)  

# Install
```
go get github.com/karldoenitz/Tigo/...
```

# Demo
## Hello Tigo
```go
package main

import (
    "github.com/karldoenitz/Tigo/TigoWeb"
    "net/http"
)

// DemoHandler handler
type DemoHandler struct {
    TigoWeb.BaseHandler
}

func (demoHandler *DemoHandler) Get() {
    demoHandler.ResponseAsText("Hello Demo!")
}

// Authorize Middleware
func Authorize(w *http.ResponseWriter, r *http.Request) bool  {
    return true
}

// Pattern
var urls = []TigoWeb.Pattern{
    {"/demo", DemoHandler{}, []TigoWeb.Middleware{Authorize}},
}

func main() {
    application := TigoWeb.Application{
        IPAddress:   "127.0.0.1",
        Port:        8888,
        UrlPatterns: urls,
    }
    application.Run()
}
```
### Compile
Open terminal, cd to target directory, input the command：
```
go build main.go
```
### Run
After compiled, there will be a runnable file named ```main```, input the command：
```
./main
```
The info will display in terminal：
```
 INFO     2022/10/07 22:40:36  Server run on: 127.0.0.1:8080
```
Open web browser and visit ```http://127.0.0.1:8888/hello-tigo```, you will see <font color=red>Hello Tigo</font>.

# Performance Comparison
<img src="https://raw.githubusercontent.com/karldoenitz/Tigo/master/documentation/chart.png" width="100%" height="300px" alt="Performance Comparison"/>

# Documentation
[Click Here](https://github.com/karldoenitz/Tigo/blob/master/documentation/documentation_en.md)

# Users of Tigo
<table>
<tr>
<td><a href="https://www.cubebackup.com" target="_blank"><img src="https://karldoenitz.github.io/TigoOld/img/cubebackup.png" width="150px" height="150px" alt="CubeBackup"/></a></td>
<td><a href="https://open2.campus.qq.com/v2/#/index/sp" target="_blank"><img src="https://karldoenitz.github.io/TigoOld/img/tencent.png" width="150px" height="150px" alt="Tencent"/></a></td>
<td><img src="https://karldoenitz.github.io/TigoOld/img/xiaomi.png" width="150px" height="150px" alt="Xiaomi"/></td>
</tr>
</table>

# Special Thanks
<table>
<tr>
<td><a href="https://www.jetbrains.com/?from=Tigo" target="_blank"><img src="https://karldoenitz.github.io/TigoOld/img/jetbrains.png" width="150px" height="150px" alt="Jetbrains"/></a></td>
</tr>
</table>

# Attention
If you like the framework, join us please.
