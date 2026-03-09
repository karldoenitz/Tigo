package web

import (
	"net/http"
	"reflect"

	"github.com/gorilla/websocket"
	"github.com/karldoenitz/Tigo/logger"
)

// 定义一个 Upgrader，用于处理 HTTP 升级请求
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 开发时允许所有跨域，生产环境应严格校验 Origin 头
	},
}

// SetWSUpgrader 设置 websocket 升级器
//   - u: 自定义的升级器
func SetWSUpgrader(u *websocket.Upgrader) {
	upgrader = *u
}

// newWSPatternHandle 创建新的 WebSocket 模式处理器
//   - handlerType: WebSocket 处理器的反射类型
//   - pattern: URL 模式配置，包含请求路径信息
//   - 返回值: 初始化完成的 WSPatternHandle 结构体指针
func newWSPatternHandle(handlerType reflect.Type, pattern Pattern) *WSPatternHandle {
	return &WSPatternHandle{
		Handler:    handlerType,
		requestUrl: pattern.Url,
	}
}

type WSPatternHandle struct {
	Handler    reflect.Type
	requestUrl string
}

// Handle 处理WebSocket连接请求
// 功能:
//
//	将HTTP连接升级为WebSocket连接，加载对应的处理器并执行通信逻辑
//	如果升级失败会记录错误日志并返回
//	连接成功后自动调用处理器的Communicate方法进行通信
//
//	- responseWriter: HTTP响应写入器
//	- request: HTTP请求对象
func (wsHandler WSPatternHandle) Handle(responseWriter http.ResponseWriter, request *http.Request) {
	// 加载handler
	handler := reflect.New(wsHandler.Handler)
	// 1. 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		logger.Error.Println("升级失败:", err)
		return
	}
	defer func() {
		// 关闭连接
		_ = conn.Close()
	}()
	// 把conn给handler.conn字段，赋上值
	handler.Elem().FieldByName("conn").Set(reflect.ValueOf(conn))
	logger.Info.Printf("websocket | %s | connect success", wsHandler.requestUrl)
	// 调用 Communicate 方法
	VoidFuncCall(handler, FnWebSocketCommunicate)
}
