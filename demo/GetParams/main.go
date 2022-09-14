package main

import (
	"fmt"
	"github.com/karldoenitz/Tigo/TigoWeb"
)

// TestHandler Demo
type TestHandler struct {
	TigoWeb.BaseHandler
}

// Post Http
func (testHandler *TestHandler) Post() {
	params := &struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Gender int    `json:"gender"`
	}{}
	paramOne := testHandler.GetParameter("one").ToBool(false)
	paramTwo := testHandler.GetParameter("two").ToFloat64()
	testHandler.GetParameter("info").To(params)
	fmt.Println(params.Name)
	fmt.Println(params.Age)
	fmt.Println(params.Gender)
	fmt.Println(paramTwo)
	if paramOne {
		testHandler.Response(paramTwo)
	}
}

func (testHandler *TestHandler) Get() {
	userId, _ := testHandler.GetPathParam("userId").ToInt()
	result := map[string]interface{}{
		"id":   userId,
		"name": "Three Zhang",
		"desc": "this is test info",
	}
	isShowExtra := testHandler.GetPathParam("isShowExtra").ToBool()
	if isShowExtra {
		result["extra"] = "this is extra info"
	}
	testHandler.ResponseAsJson(result)
}

var url = []TigoWeb.Pattern{
	{"/test", TestHandler{}, nil},
	{"/user/{userId}/info/{isShowExtra}", TestHandler{}, nil},
}

func main() {
	application := TigoWeb.Application{
		IPAddress:   "0.0.0.0",
		Port:        8080,
		UrlPatterns: url,
	}
	application.Run()
}
