// tiger插件，一个脚手架工具，用于来初始化一个Tigo项目
package main

const DemoCode = `
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

func main() {
	// TODO 此处补充脚手架逻辑
	// 获取命令行参数，根据参数判断是否是创建demo，
	// 如果创建demo，则直接把`DemoCode`注入到目标文件中就行
	// tiger支持的命令:
	//  - create xxx: 创建项目
	//  - addHandler xxx: 增加xxx命名的handler
	//  - conf xxx: 用xxx命名的配置文件替换现有配置文件，没有则新建
}
