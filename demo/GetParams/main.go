package main

import (
	"github.com/karldoenitz/Tigo/TigoWeb"
	"fmt"
)

type TestHandler struct {
	TigoWeb.BaseHandler
}

func (testHandler *TestHandler) Post() {
	params := &struct {
		Name string `json:"name"`
		Age int `json:"age"`
		Gender int `json:"gender"`
	}{}
	paramOne := testHandler.GetParameter("one").ToBool(false)
	paramTwo := testHandler.GetParameter("two")
	testHandler.GetParameter("info").To(params)
	fmt.Println(params.Name)
	fmt.Println(params.Age)
	fmt.Println(params.Gender)
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
