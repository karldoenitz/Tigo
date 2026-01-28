package web

import (
	"net/http"
	"reflect"
)

// UrlPatternHandle 是URL路由句柄，用来驱动url路由以及其映射的handler
type UrlPatternHandle struct {
	Handler    reflect.Type
	requestUrl string
}

// Handle 封装HTTP请求的中间件，主要有以下功能：
//   - 1、根据反射找到挂载的handler；
//   - 2、调用handler的InitHandler方法；
//   - 3、进行HTTP请求预处理，包括判断请求方式是否合法等；
//   - 4、调用handler中的功能方法；
//   - 5、进行HTTP请求结束处理。
func (urlPatternMidWare UrlPatternHandle) Handle(responseWriter http.ResponseWriter, request *http.Request) {
	// 加载handler
	handler := reflect.New(urlPatternMidWare.Handler)
	// 调用InitHandler方法
	VoidFuncCall(handler, FnInitHandler, reflect.ValueOf(responseWriter), reflect.ValueOf(request))
	// 调用PassJson方法
	VoidFuncCall(handler, FnPassJson)
	// 调用BeforeRequest方法
	VoidFuncCall(handler, FnBeforeRequest)
	// 根据http请求方式调用相关方法
	VoidFuncCall(handler, MethodEnum(request.Method))
	// 调用TeardownRequest方法
	VoidFuncCall(handler, FnTeardownRequest)
}

// newUrlPatternHandle 创建并返回一个新的UrlPatternHandle实例
//   - handlerType: 处理器的反射类型
//   - pattern: 包含URL模式的路由配置
func newUrlPatternHandle(handlerType reflect.Type, pattern Pattern) *UrlPatternHandle {
	return &UrlPatternHandle{
		Handler:    handlerType,
		requestUrl: pattern.Url,
	}
}
