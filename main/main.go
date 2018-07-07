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
	value := helloHandler.GetParameter("hello")
	fmt.Println(value)
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

type RedirectHandler struct {
	TigoWeb.BaseHandler
}

func (redirectHandler *RedirectHandler)Handle(responseWriter http.ResponseWriter, request *http.Request) {
	redirectHandler.InitHandler(responseWriter, request)
	if !redirectHandler.CheckRequestMethodValid("GET") {
		return
	}
	redirectHandler.Redirect("http://www.baidu.com")
}

var urls = map[string]interface{Handle(http.ResponseWriter, *http.Request)}{
	"/hello-world": &HelloHandler{},
	"/redirect"   : &RedirectHandler{},
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
