package main

import (
	"github.com/karldoenitz/Tigo/TigoWeb"
	"fmt"
	"github.com/karldoenitz/Tigo/logger"
)

type HelloHandler struct {
	TigoWeb.BaseHandler
}

func (helloHandler *HelloHandler)Handle() {
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
	c, _ := helloHandler.Request.Cookie("Tigo")
	fmt.Println(c)
	helloHandler.ResponseAsHtml("<p1 style='color: red'>Hello Go!</p1>")
}

type RedirectHandler struct {
	TigoWeb.BaseHandler
}

func (redirectHandler *RedirectHandler)Handle() {
	if !redirectHandler.CheckRequestMethodValid("GET") {
		return
	}
	redirectHandler.Redirect("http://www.tencent.com")
}

type TestCookieHandler struct {
	TigoWeb.BaseHandler
}

func (testCookieHandler *TestCookieHandler)Handle() {
	testCookieHandler.SetSecureCookie("name", "value")
	cookie := testCookieHandler.GetSecureCookie("name")
	fmt.Println(cookie)
	logger.Trace.Printf("no trace: %s", "info")
	logger.Info.Printf("this is info: %s", "Info Here")
	logger.Warning.Printf("this is warning info: %s", "Warning Here")
	logger.Error.Printf("interesting %s", "testing")
	testCookieHandler.ResponseAsHtml("<h1>Tiger Go Go Go!</h1>")
}

var urls = map[string]interface{}{
	"/hello-world": &HelloHandler{},
	"/redirect"   : &RedirectHandler{},
	"/test-cookie": &TestCookieHandler{},
}

func main() {
	application := TigoWeb.Application{
		IPAddress:  "0.0.0.0",
		Port:       8888,
		UrlPattern: urls,
		ConfigPath: "./configuration.json",
	}
	application.Run()
}
