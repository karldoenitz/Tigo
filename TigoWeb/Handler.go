// Copyright 2018 The Tigo Authors. All rights reserved.
package WebFramework

import (
	"net/http"
	"encoding/json"
	"fmt"
	"strings"
)

// Handler的基础类，开发者开发的handler继承此类
type BaseHandler struct {
	ResponseWriter  http.ResponseWriter
	Request        *http.Request
}

// 初始化Handler的方法
func (baseHandler *BaseHandler)InitHandler(responseWriter http.ResponseWriter, request *http.Request) {
	baseHandler.Request = request
	baseHandler.ResponseWriter = responseWriter
	baseHandler.Request.ParseForm()
}

// 将对象转化为Json字符串，转换失败则返回空字符串。
// 传入参数Response为一个interface，必须有成员函数Print。
func (baseHandler *BaseHandler)ToJson(response Response) (result string) {
	// 将该对象转换为byte字节数组
	jsonResult, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		return ""
	}
	// 将byte数组转换为string
	return string(jsonResult)
}

// 向客户端响应一个Json结果
func (baseHandler *BaseHandler)ResponseAsJson(response Response)  {
	// 将对象转换为Json字符串
	jsonResult := baseHandler.ToJson(response)
	// 设置http报文头内的Content-Type
	baseHandler.ResponseWriter.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(baseHandler.ResponseWriter, jsonResult)
}

// 向客户端响应一个Text结果
func (baseHandler *BaseHandler)ResponseAsText(result string)  {
	fmt.Fprintf(baseHandler.ResponseWriter, result)
}

// 向客户端响应一个html结果
func (baseHandler *BaseHandler)ResponseAsHtml(result string)  {
	baseHandler.ResponseWriter.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(baseHandler.ResponseWriter, result)
}

// 检查请求是否被允许
func (baseHandler *BaseHandler)CheckRequestMethodValid(methods ...string)(result bool) {
	// 获取请求方式
	requestMethod := baseHandler.Request.Method
	// 遍历被允许的请求方式，判断是否合法
	for _, value := range methods {
		if requestMethod == value || strings.ToLower(requestMethod) == value {
			return true
		}
	}
	// 如果不合法返回405
	baseHandler.ResponseWriter.WriteHeader(405)
	return false
}

// 设置cookie
func (baseHandler *BaseHandler)SetCookie(name string, value string) {

}

// 设置高级cookie选项
func (baseHandler *BaseHandler)SetCookieObject(cookie http.Cookie) {

}

// 设置加密cookie

// 获取cookie值

// 获取加密cookie值

// 获取cookie对象

// 获取header

// 设置header
