[![Build Status](https://travis-ci.org/karldoenitz/Tigo.svg?branch=master)](https://travis-ci.org/karldoenitz/Tigo)
# Tigo([For English Documentation Click Here](https://github.com/karldoenitz/Tigo/blob/master/README_EN.md))
![Tigo logo](https://github.com/karldoenitz/karlooper/blob/master/documentations/images/logo.jpg "this is Tigo logo")  
一个使用Go语言开发的web框架。

# 安装
```
go get github.com/karldoenitz/Tigo/TigoWeb
```

# 示例
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
var urls = map[string]interface{Handle(http.ResponseWriter, *http.Request)}{
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
### 编译
打开终端，进入代码目录，运行如下命令：
```
go build main.go
```
### 运行
编译完成后，会有一个可执行文件```main```，运行如下命令：
```
./main
```
终端会有如下显示：
```
INFO: 2018/07/09 15:02:36 Application.go:22: Server run on: 0.0.0.0:8888
```
打开浏览器访问地址```http://127.0.0.1:8888/hello-tigo```，就可以看到<font color=red>Hello Tigo</font>。

# 文档
[点击此处](https://github.com/karldoenitz/Tigo/blob/master/documentation/documentation.md)

# 注意
这个框架在Linux版本的 [CubeBackup for Google Apps](http://www.cubebackup.com) 中有所使用。  
如果你对次框架感兴趣，可以加入我们一同开发。
