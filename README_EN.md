[![Badge](https://img.shields.io/badge/link-Tigo-blue.svg)](https://karldoenitz.github.io/Tigo-EN/)
[![LICENSE](https://img.shields.io/badge/license-Tigo-blue.svg)](https://github.com/karldoenitz/Tigo/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/karldoenitz/Tigo.svg?branch=master)](https://travis-ci.org/karldoenitz/Tigo)
[![Join the chat at https://gitter.im/karlooper/Tigo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/karlooper/Tigo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Open Source Helpers](https://www.codetriage.com/karldoenitz/tigo/badges/users.svg)](https://www.codetriage.com/karldoenitz/Tigo)
[![Go Report Card](https://goreportcard.com/badge/github.com/karldoenitz/Tigo)](https://goreportcard.com/report/github.com/karldoenitz/Tigo)
[![GoDoc](https://godoc.org/github.com/karldoenitz/Tigo?status.svg)](https://godoc.org/github.com/karldoenitz/Tigo)
[![Release](https://img.shields.io/github/release/karldoenitz/Tigo.svg)](https://github.com/karldoenitz/Tigo/releases)  
![Tigo logo](https://github.com/karldoenitz/Tigo/blob/master/documentation/tigo_logo.jpg "this is Tigo logo")
# Tigo([中文文档点击此处](https://github.com/karldoenitz/Tigo/blob/master/README.md))
A web framework developed in go language.

# Recommend A Commandline Tool For U!
`tiger` is a commandline tool for `Tigo` framework, you can use `tiger` to create a `Tigo` projection.  
[glance tiger](https://github.com/karldoenitz/tiger)

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

// handler
type DemoHandler struct {
    TigoWeb.BaseHandler
}

func (demoHandler *DemoHandler) Get() {
    demoHandler.ResponseAsText("Hello Demo!")
}

// Middleware
func Authorize(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 此处授权认证逻辑
        next.ServeHTTP(w, r)
    }
}

// Router
var urls = []TigoWeb.Router{
    {"/demo", &DemoHandler{}, []TigoWeb.Middleware{Authorize}},
}

func main() {
    application := TigoWeb.Application{
        IPAddress:  "127.0.0.1",
        Port:       8888,
        UrlRouters: urls,
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
INFO: 2018/07/09 15:02:36 Application.go:22: Server run on: 127.0.0.1:8888
```
Open web browser and visit ```http://127.0.0.1:8888/hello-tigo```, you will see <font color=red>Hello Tigo</font>.

# Documentation
[Click Here](https://github.com/karldoenitz/Tigo/blob/master/documentation/documentation_en.md)

# Users of Tigo
<table>
<tr>
<td><a href="https://www.cubebackup.com" target="_blank"><img src="https://karldoenitz.github.io/TigoOld/img/cubebackup.png" width="150px" height="150px"/></a></td>
<td><img src="https://karldoenitz.github.io/TigoOld/img/tencent.png" width="150px" height="150px"/></td>
<td><img src="https://karldoenitz.github.io/TigoOld/img/xiaomi.png" width="150px" height="150px"/></td>
</tr>
</table>

# Attention
If you like the framework, join us please.
