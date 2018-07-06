package main

import (
	"net/http"
	"../TigoWeb"
	"fmt"
)

type HelloHandler struct {
	TigoWeb.BaseHandler
}

func (helloHandler *HelloHandler)Handle(responseWriter http.ResponseWriter, request *http.Request) {
	helloHandler.InitHandler(responseWriter, request)
	if !helloHandler.CheckRequestMethodValid("GET", "POST") {
		return
	}
	cookie := TigoWeb.Cookie{
		Name:        "Tigo",
		Value:       "Tiger Go",
		SecurityKey: "player",
		Path:        "/hello-world",
		IsSecurity:  true,
	}
	helloHandler.SetCookieObject(cookie)
	clientCookie, _ := helloHandler.GetCookieObject("Tigo", "player")
	fmt.Println(clientCookie.Name)
	fmt.Println(clientCookie.Value)
	fmt.Println(clientCookie.SecurityKey)
	fmt.Println(clientCookie.Path)
	c, _ := request.Cookie("Tigo")
	fmt.Println(c)
	helloHandler.ResponseAsHtml("<p1 style='color: red'>Hello Go!</p1>")
}

var urls = map[string]interface{Handle(http.ResponseWriter, *http.Request)}{
	"/hello-world": &HelloHandler{},
}

func main() {
	urlPattern := TigoWeb.UrlPattern{UrlMapping: urls}
	application := TigoWeb.Application{
		IPAddress:  "0.0.0.0",
		Port:       "8888",
		UrlPattern: urlPattern,
	}
	application.Run()
}
