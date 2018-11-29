// Copyright 2018 The Tigo Authors. All rights reserved.
package TigoWeb

import (
	"github.com/karldoenitz/Tigo/logger"
	"net/http"
	"reflect"
	"time"
)

// URL路由中间件
type UrlPatternMidWare struct {
	Handler    interface{}
	requestUrl string
}

// 封装HTTP请求的中间件，主要有以下功能：
//  - 1、根据反射找到挂载的handler；
//  - 2、调用handler的InitHandler方法；
//  - 3、进行HTTP请求预处理，包括判断请求方式是否合法等；
//  - 4、调用handler中的功能方法；
//  - 5、进行HTTP请求结束处理。
func (urlPatternMidWare UrlPatternMidWare) Handle(responseWriter http.ResponseWriter, request *http.Request) {
	requestStart := time.Now().Nanosecond() / 1e6
	// 加载handler
	handler := reflect.ValueOf(urlPatternMidWare.Handler)
	// 获取init方法
	init := handler.MethodByName("InitHandler")
	// 解析参数
	paramPasser := handler.MethodByName("PassJson")
	// 获取HTTP请求方式
	requestMethod := MethodMapping[request.Method]
	function := handler.MethodByName(requestMethod)
	initParams := []reflect.Value{reflect.ValueOf(responseWriter), reflect.ValueOf(request)}
	var functionParams []reflect.Value
	if init.IsValid() {
		init.Call(initParams)
	}
	if paramPasser.IsValid() {
		paramPasser.Call(functionParams)
	}
	if function.IsValid() {
		function.Call(functionParams)
	}
	requestEnd := time.Now().Nanosecond() / 1e6
	logger.Trace.Printf("%s %s %dms", requestMethod, urlPatternMidWare.requestUrl, requestEnd-requestStart)
}

// URL路由，此处存储URL映射。
type UrlPattern struct {
	UrlMapping map[string]interface{}
}

// 向http服务挂载单个handler，注意：
//   - handler必须有一个Handle(http.ResponseWriter, *http.Request)函数
func (urlPattern *UrlPattern) AppendUrlPattern(uri string, v interface {
	Handle(http.ResponseWriter, *http.Request)
}) {
	http.HandleFunc(uri, v.Handle)
}

// 初始化url映射，遍历UrlMapping，将handler与对应的URL依次挂载到http服务上
func (urlPattern *UrlPattern) Init() {
	for key, value := range urlPattern.UrlMapping {
		urlPatternMidWare := UrlPatternMidWare{
			Handler:    value,
			requestUrl: key,
		}
		urlPattern.AppendUrlPattern(key, &urlPatternMidWare)
	}
}
