package main

import "github.com/karldoenitz/Tigo/TigoWeb"

type TestHandler struct {
	TigoWeb.BaseHandler
}

func (testHandler *TestHandler) Post() {
	paramOne := testHandler.GetParameter("one").(bool)
	paramTwo := testHandler.GetParameter("two")
	if paramOne {
		testHandler.Response(paramTwo)
	}
}

var url = map[string] interface{} {
	"/test": &TestHandler{},
}

func main() {
	application := TigoWeb.Application{
		IPAddress:"0.0.0.0",
		Port:8080,
		UrlPattern:url,
	}
	application.Run()
}
