[![Build Status](https://travis-ci.org/karldoenitz/Tigo.svg?branch=master)](https://travis-ci.org/karldoenitz/Tigo)  
![Tigo logo](https://github.com/karldoenitz/Tigo/blob/master/documentation/tigo_logo.jpg "this is Tigo logo")
# Tigo([中文文档点击此处](https://github.com/karldoenitz/Tigo/blob/master/README.md))
A web framework developed in go language.

# Install
```
go get github.com/karldoenitz/Tigo/...
```

# Demo
## Hello Tigo
```go
package main

import (
    "net/http"
    "github.com/karldoenitz/Tigo/TigoWeb"
)

// handler
type HelloHandler struct {
    TigoWeb.BaseHandler
}

func (helloHandler *HelloHandler)Handle(responseWriter http.ResponseWriter, request *http.Request) {
    helloHandler.InitHandler(responseWriter, request)
    helloHandler.ResponseAsHtml("<p1 style='color: red'>Hello Tigo!</p1>")
}

// url路由配置
var urls = map[string]interface{}{
    "/hello-tigo": &HelloHandler{},
}

// 主函数
func main() {
    urlPattern := TigoWeb.UrlPattern{UrlMapping: urls}
    application := TigoWeb.Application{
        IPAddress:  "0.0.0.0",
        Port:       8888,
        UrlPattern: urlPattern,
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
INFO: 2018/07/09 15:02:36 Application.go:22: Server run on: 0.0.0.0:8888
```
Open web browser and visit ```http://127.0.0.1:8888/hello-tigo```, you will see <font color=red>Hello Tigo</font>.

# Documentation
[Click Here](https://github.com/karldoenitz/Tigo/blob/master/documentation/documentation_en.md)

# Attention
This framework used in [CubeBackup for Google Apps](http://www.cubebackup.com)。
If you like the framework, join us please.
