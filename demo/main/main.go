package main

import (
	"fmt"
	"github.com/karldoenitz/Tigo/TigoWeb"
	"github.com/karldoenitz/Tigo/logger"
)

// Demo HelloHandler
type HelloHandler struct {
	TigoWeb.BaseHandler
}

// Http Get
func (helloHandler *HelloHandler) Get() {
	helloHandler.DumpHttpRequestMsg(logger.TraceLevel)
	value := helloHandler.GetParameter("hello").ToString()
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

// Demo RedirectHandler
type RedirectHandler struct {
	TigoWeb.BaseHandler
}

// Http Get
func (redirectHandler *RedirectHandler) Get() {
	redirectHandler.Redirect("http://www.tencent.com")
}

// Demo TestCookieHandler
type TestCookieHandler struct {
	TigoWeb.BaseHandler
}

// Http Get
func (testCookieHandler *TestCookieHandler) Get() {
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
	"/redirect":    &RedirectHandler{},
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
