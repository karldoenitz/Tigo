[![Badge](https://img.shields.io/badge/link-Tigo-blue.svg)](https://karldoenitz.github.io/Tigo/)
[![LICENSE](https://img.shields.io/badge/license-Tigo-blue.svg)](https://github.com/karldoenitz/Tigo/blob/master/LICENSE)
[![Go](https://github.com/karldoenitz/Tigo/actions/workflows/go.yml/badge.svg)](https://github.com/karldoenitz/Tigo/actions/workflows/go.yml)
[![Join the chat at https://gitter.im/karlooper/Tigo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/karlooper/Tigo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Open Source Helpers](https://www.codetriage.com/karldoenitz/tigo/badges/users.svg)](https://www.codetriage.com/karldoenitz/Tigo)
[![Go Report Card](https://goreportcard.com/badge/github.com/karldoenitz/Tigo)](https://goreportcard.com/report/github.com/karldoenitz/Tigo)
[![GoDoc](https://godoc.org/github.com/karldoenitz/Tigo?status.svg)](https://pkg.go.dev/github.com/karldoenitz/Tigo)
[![Release](https://img.shields.io/github/release/karldoenitz/Tigo.svg)](https://github.com/karldoenitz/Tigo/releases)  
![Tigo logo](https://raw.githubusercontent.com/karldoenitz/Tigo/master/documentation/tigo_logo.jpg "this is Tigo logo")  
# Tigo([For English Documentation Click Here](https://github.com/karldoenitz/Tigo/blob/master/README_EN.md))
一个使用Go语言开发的web框架。

# 相关工具及插件
- **tiger**  
`tiger`是一个专门为`Tigo`框架量身定做的脚手架工具，可以使用`tiger`新建`Tigo`项目或者执行其他操作。  
[查看tiger](https://github.com/karldoenitz/tiger)  
- **tission**  
`tission`是一个为`Tigo`定制的session插件。  
[查看tission](https://github.com/karldoenitz/tission)

# 安装
```shell
export GO111MODULE=off; 
go get github.com/karldoenitz/Tigo/...
export GO111MODULE=on; 
```

# 升级
```shell
export GO111MODULE=off; 
go get -u github.com/karldoenitz/Tigo/...
export GO111MODULE=on; 
```

# 示例
## Hello Tigo

```go
package main

import (
    "github.com/karldoenitz/Tigo/web"
    "net/http"
)

// DemoHandler handler
type DemoHandler struct {
    web.BaseHandler
}

func (demoHandler *DemoHandler) Get() {
    demoHandler.ResponseAsText("Hello Demo!")
}

// Authorize 中间件
func Authorize(w *http.ResponseWriter, r *http.Request) bool {
    // 此处返回true表示继续执行，false则直接返回，后续的中间件不会执行 
    return true
}

// 路由
var urls = []web.Pattern{
    {"/demo", DemoHandler{}, []web.Middleware{Authorize}},
}

func main() {
    application := web.Application{
        IPAddress:   "127.0.0.1",
        Port:        8888,
        UrlPatterns: urls,
    }
    application.Run()
}
```
### 编译
打开终端，进入代码目录，运行如下命令：
```shell
go build main.go
```
### 运行
编译完成后，会有一个可执行文件```main```，运行如下命令：
```shell
./main
```
终端会有如下显示：
```
 INFO     2022/10/07 22:40:36  Server run on: http://127.0.0.1:8080
```
打开浏览器访问地址[```http://127.0.0.1:8888/demo```](http://127.0.0.1:8888/demo)，就可以看到<font color=red>Hello Demo</font>。

# 性能对比
<img src="https://raw.githubusercontent.com/karldoenitz/Tigo/master/documentation/chart.png" width="100%" height="300px" alt="性能对比"/> 

# 文档
[点击此处](https://github.com/karldoenitz/Tigo/blob/master/documentation/documentation.md)
# 都有谁在使用Tigo
<table>
<tr>
<td><a href="https://www.cubebackup.com" target="_blank"><img src="https://karldoenitz.github.io/TigoOld/img/cubebackup.png" width="150px" height="150px" alt="cube-backup"/></a></td>
<td><a href="https://open2.campus.qq.com" target="_blank"><img src="https://karldoenitz.github.io/TigoOld/img/tencent.png" width="150px" height="150px" alt="Tencent"/></a></td>
<td><img src="https://karldoenitz.github.io/TigoOld/img/xiaomi.png" width="150px" height="150px" alt="Xiaomi"/></td>
</tr>
</table>

# 鸣谢以下组织的支持
<table>
<tr>
<td><a href="https://www.jetbrains.com/?from=Tigo" target="_blank"><img src="https://karldoenitz.github.io/TigoOld/img/jetbrains.png" width="150px" height="150px" alt="Jetbrains"/></a></td>
</tr>
</table>


# 注意
如果你对此框架感兴趣，可以加入我们一同开发。
