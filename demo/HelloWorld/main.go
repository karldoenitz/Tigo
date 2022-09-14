package main

import "github.com/karldoenitz/Tigo/TigoWeb"

// HelloWorldHandler Demo
type HelloWorldHandler struct {
	TigoWeb.BaseHandler
}

// Get Http Method
func (helloWorldHandler *HelloWorldHandler) Get() {
	helloWorldHandler.ResponseAsHtml("Hello World!")
}

var urls = []TigoWeb.Pattern{
	{"/hello-world", HelloWorldHandler{}, nil},
}

func main() {
	application := TigoWeb.Application{
		IPAddress:   "0.0.0.0",
		Port:        8080,
		UrlPatterns: urls,
	}
	application.Run()
}
