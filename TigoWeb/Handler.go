// Copyright 2018 The Tigo Authors. All rights reserved.
package WebFramework

import (
	"net/http"
	"encoding/json"
	"fmt"
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
