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

var urls = map[string]interface{}{
	"/hello-world": &HelloWorldHandler{},
}

func main() {
	application := TigoWeb.Application{
		IPAddress:  "0.0.0.0",
		Port:       8080,
		UrlPattern: urls,
	}
	application.Run()
}
