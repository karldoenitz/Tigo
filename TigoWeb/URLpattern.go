// Package TigoWeb Copyright 2018 The Tigo Authors. All rights reserved.
package TigoWeb

import (
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
	// 获取init方法
	init := handler.MethodByName("InitHandler")
	// 解析参数
	paramPasser := handler.MethodByName("PassJson")
	// 获取BeforeRequest方法
	beforeRequest := handler.MethodByName("BeforeRequest")
	// 获取HTTP请求方式
	requestMethod := MethodMapping[request.Method]
	function := handler.MethodByName(requestMethod)
	// 获取TeardownRequest方法
	teardownRequest := handler.MethodByName("TeardownRequest")
	initParams := []reflect.Value{reflect.ValueOf(responseWriter), reflect.ValueOf(request)}
	var functionParams []reflect.Value
	if init.IsValid() {
		init.Call(initParams)
	}
	if paramPasser.IsValid() {
		paramPasser.Call(functionParams)
	}
	if beforeRequest.IsValid() {
		beforeRequest.Call(functionParams)
	}
	if function.IsValid() {
		function.Call(functionParams)
	}
	if teardownRequest.IsValid() {
		teardownRequest.Call(functionParams)
	}
}

// Router 路由对象
type Router struct {
	Url        string
	Handler    interface{}
	Middleware []Middleware
}

// UrlPattern 是URL路由，此处存储URL映射。
type UrlPattern struct {
	UrlMapping map[string]interface{}
	UrlRouters []Router
}

// AppendUrlPattern 向http服务挂载单个handler，注意：
//   - handler必须有一个Handle(http.ResponseWriter, *http.Request)函数
func (urlPattern *UrlPattern) AppendUrlPattern(uri string, v interface {
	Handle(http.ResponseWriter, *http.Request)
}) {
	http.HandleFunc(uri, v.Handle)
}

// AppendRouterPattern 向http服务挂载单个Router，Router中配置有url对应的handler以及对应的中间件
func (urlPattern *UrlPattern) AppendRouterPattern(router Router, v interface {
	Handle(http.ResponseWriter, *http.Request)
}) {
	baseMiddleware := []Middleware{HttpContextLogMiddleware, InternalServerErrorMiddleware}
	baseMiddleware = append(baseMiddleware, router.Middleware...)
	middleware := chainMiddleware(baseMiddleware...)
	http.HandleFunc(router.Url, middleware(v.Handle))
}

// Init 初始化url映射，遍历UrlMapping，将handler与对应的URL依次挂载到http服务上
func (urlPattern *UrlPattern) Init() {
	for key, value := range urlPattern.UrlMapping {
		urlPatternMidWare := UrlPatternMidWare{
			Handler:    value,
			requestUrl: key,
		}
		urlPattern.AppendUrlPattern(key, &urlPatternMidWare)
	}
	for _, router := range urlPattern.UrlRouters {
		urlPatternMidWare := UrlPatternMidWare{
			Handler:    router.Handler,
			requestUrl: router.Url,
		}
		urlPattern.AppendRouterPattern(router, &urlPatternMidWare)
	}
}
