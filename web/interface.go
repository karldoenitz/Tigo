// Package web Copyright 2018 The Tigo Authors. All rights reserved.
package web

// Response 响应给客户端的interface，用户自定义实现
type Response interface {
	Print()
}
