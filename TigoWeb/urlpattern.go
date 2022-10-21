// Package TigoWeb Copyright 2018 The Tigo Authors. All rights reserved.
package TigoWeb

import (
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

// UrlPatternMidWare 是URL路由中间件
type UrlPatternMidWare struct {
	Handler    interface{}
	requestUrl string
}

// Handle 封装HTTP请求的中间件，主要有以下功能：
//  - 1、根据反射找到挂载的handler；
//  - 2、调用handler的InitHandler方法；
//  - 3、进行HTTP请求预处理，包括判断请求方式是否合法等；
//  - 4、调用handler中的功能方法；
//  - 5、进行HTTP请求结束处理。
func (urlPatternMidWare UrlPatternMidWare) Handle(responseWriter http.ResponseWriter, request *http.Request) {
	handlerType := reflect.TypeOf(urlPatternMidWare.Handler)
	if handlerType.Kind() == reflect.Ptr {
		handlerType = handlerType.Elem()
	}
	// 加载handler
	handler := reflect.New(handlerType)
	// 调用InitHandler方法
	VoidFuncCall(handler, "InitHandler", reflect.ValueOf(responseWriter), reflect.ValueOf(request))
	// 调用PassJson方法
	VoidFuncCall(handler, "PassJson")
	// 调用BeforeRequest方法
	VoidFuncCall(handler, "BeforeRequest")
	// 根据http请求方式调用相关方法
	VoidFuncCall(handler, MethodMapping[request.Method])
	// 调用TeardownRequest方法
	VoidFuncCall(handler, "TeardownRequest")
}

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
func (urlPattern *UrlPattern) AppendRouterPattern(pattern Pattern, v interface {
	Handle(http.ResponseWriter, *http.Request)
}) {
	baseMiddleware := []Middleware{HttpContextLogMiddleware, InternalServerErrorMiddleware}
	baseMiddleware = append(baseMiddleware, pattern.Middleware...)
	middleware := chainMiddleware(baseMiddleware...)
	urlPattern.router.HandleFunc(pattern.Url, middleware(v.Handle))
}

// Init 初始化url映射，遍历UrlMapping，将handler与对应的URL依次挂载到http服务上
func (urlPattern *UrlPattern) Init() {
	for _, pattern := range urlPattern.UrlPatterns {
		urlPattern.AppendRouterPattern(pattern, &UrlPatternMidWare{
			Handler:    pattern.Handler,
			requestUrl: pattern.Url,
		})
	}
}
