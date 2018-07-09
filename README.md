# Tigo
![Tigo logo](https://github.com/karldoenitz/karlooper/blob/master/documentations/images/logo.jpg "this is Tigo logo")  
一个使用Go语言开发的web框架。

# 安装
```
go get github.com/karldoenitz/Tigo
```

# 示例
## Hello Tigo
```go
package main

import (
    "net/http"
    "github.com/Tigo/TigoWeb"
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
        Port:       "8888",
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
打开浏览器访问地址```http://127.0.0.1:8888/hello-tigo```，就可以看到<span style='color: red'>Hello Tigo<span>。

# 注意
这个框架在Linux版本的 [CubeBackup for Google Apps](http://www.cubebackup.com) 中有所使用。
如果你对代码的改进有好的建议或意见,欢迎发送邮件至 karlvorndoenitz@gmail.com, 非常感谢。