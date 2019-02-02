[![Build Status](https://travis-ci.org/karldoenitz/Tigo.svg?branch=master)](https://travis-ci.org/karldoenitz/Tigo)
[![Join the chat at https://gitter.im/karlooper/Tigo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/karlooper/Tigo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Open Source Helpers](https://www.codetriage.com/karldoenitz/tigo/badges/users.svg)](https://www.codetriage.com/karldoenitz/Tigo)
[![Go Report Card](https://goreportcard.com/badge/github.com/karldoenitz/Tigo)](https://goreportcard.com/report/github.com/karldoenitz/Tigo)
[![GoDoc](https://godoc.org/github.com/karldoenitz/Tigo?status.svg)](https://godoc.org/github.com/karldoenitz/Tigo)
[![Release](https://img.shields.io/github/release/karldoenitz/Tigo.svg?style=flat-square)](https://github.com/karldoenitz/Tigo/releases)  
![Tigo logo](https://github.com/karldoenitz/Tigo/blob/master/documentation/tigo_logo.jpg "this is Tigo logo")  
# Tigo([For English Documentation Click Here](https://github.com/karldoenitz/Tigo/blob/master/README_EN.md))
一个使用Go语言开发的web框架。

# 安装
```
go get github.com/karldoenitz/Tigo/...
```

# 示例
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

// 中间件
func Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 此处授权认证逻辑
		next.ServeHTTP(w, r)
	}
}

// 路由
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
如果你对此框架感兴趣，可以加入我们一同开发。
