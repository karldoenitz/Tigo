// Package web Copyright 2018 The Tigo Authors. All rights reserved.
package web

import "net/http"

// Response 响应给客户端的interface，用户自定义实现
type Response interface {
	Print()
}

// OriginHandlerInterface 是原始的handler接口，区别于Tigo框架内的handler接口
type OriginHandlerInterface interface {
	Handle(http.ResponseWriter, *http.Request)
}
