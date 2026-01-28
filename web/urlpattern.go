// Package web Copyright 2018 The Tigo Authors. All rights reserved.
package web

import (
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

// Pattern 路由对象
type Pattern struct {
	Url        string
	Handler    interface{}
	Middleware []Middleware
}

// UrlPattern 是URL路由，此处存储URL映射。
type UrlPattern struct {
	UrlPatterns []Pattern
	router      *mux.Router
}

// AppendRouterPattern 向http服务挂载单个Router，Router中配置有url对应的handler以及对应的中间件
func (urlPattern *UrlPattern) AppendRouterPattern(pattern Pattern, v OriginHandlerInterface) {
	// 判断是否是文件服务器
	if filePath, isFileServer := pattern.Handler.(string); isFileServer {
		fileServer := http.FileServer(http.Dir(filePath))
		fileRouter := urlPattern.router.PathPrefix(pattern.Url).Subrouter()
		var fileServerMiddleWares []mux.MiddlewareFunc
		for _, mv := range pattern.Middleware {
			m := convertHandleMV(mv)
			fileServerMiddleWares = append(fileServerMiddleWares, m)
		}
		fileRouter.Use(fileServerMiddleWares...)
		fileRouter.PathPrefix("/").Handler(http.StripPrefix(pattern.Url, fileServer))
		return
	}
	// 判断是否是Websocket
	if _, isWSHandler := v.(*WSPatternHandle); isWSHandler {
		var wsMiddleware []middleware
		for _, mv := range pattern.Middleware {
			m := convertHandleFuncMV(mv)
			wsMiddleware = append(wsMiddleware, m)
		}
		wsMiddlewares := chainMiddleware(wsMiddleware...)
		urlPattern.router.HandleFunc(pattern.Url, wsMiddlewares(v.Handle))
		return
	}
	// 判断是否是handler
	baseMiddleware := []middleware{HttpContextLogMiddleware, InternalServerErrorMiddleware}
	for _, mv := range pattern.Middleware {
		m := convertHandleFuncMV(mv)
		baseMiddleware = append(baseMiddleware, m)
	}
	middlewares := chainMiddleware(baseMiddleware...)
	urlPattern.router.HandleFunc(pattern.Url, middlewares(v.Handle))
}

// Init 初始化url映射，遍历UrlMapping，将handler与对应的URL依次挂载到http服务上
func (urlPattern *UrlPattern) Init() {
	for _, pattern := range urlPattern.UrlPatterns {
		oriHandler := NewOriginHandler(pattern)
		urlPattern.AppendRouterPattern(pattern, oriHandler)
	}
}

func NewOriginHandler(pattern Pattern) OriginHandlerInterface {
	handlerType := reflect.TypeOf(pattern.Handler)
	if handlerType.Kind() == reflect.Ptr {
		handlerType = handlerType.Elem()
	}
	_, isWSHandler := handlerType.FieldByName(FnWebSocketBaseHandler)
	if isWSHandler {
		return newWSPatternHandle(handlerType, pattern)
	}
	return newUrlPatternHandle(handlerType, pattern)
}
