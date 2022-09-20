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

}
