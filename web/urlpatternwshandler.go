package web

import (
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
	"github.com/karldoenitz/Tigo/logger"
)

// 全局 WebSocket Upgrader
var globalWSUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	HandshakeTimeout: 10 * time.Second,
}

// SetGlobalWSUpgrader 设置全局 WebSocket 升级器
// 使用此方法可以统一配置所有 WebSocket 连接的升级行为
//   - u: 自定义的升级器
func SetGlobalWSUpgrader(u *websocket.Upgrader) {
	globalWSUpgrader = *u
}

// GetGlobalWSUpgrader 获取全局 WebSocket 升级器
func GetGlobalWSUpgrader() *websocket.Upgrader {
	return &globalWSUpgrader
}

// WSHandlerConfig WebSocket handler 配置
type WSHandlerConfig struct {
	CheckOrigin       func(*http.Request) bool
	HandshakeTimeout  time.Duration
	ReadBufferSize    int
	WriteBufferSize   int
	EnableCompression bool
}

// DefaultWSHandlerConfig 返回默认的 WebSocket handler 配置
func DefaultWSHandlerConfig() WSHandlerConfig {
	return WSHandlerConfig{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		HandshakeTimeout:  10 * time.Second,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: false,
	}
}

// newWSPatternHandle 创建新的 WebSocket 模式处理器
//   - handlerType: WebSocket 处理器的反射类型
//   - pattern: URL 模式配置，包含请求路径信息
//   - 返回值: 初始化完成的 WSPatternHandle 结构体指针
func newWSPatternHandle(handlerType reflect.Type, pattern Pattern) *WSPatternHandle {
	return &WSPatternHandle{
		Handler:    handlerType,
		requestUrl: pattern.Url,
		config:     DefaultWSHandlerConfig(),
	}
}

// WSPatternHandle WebSocket 模式处理器
type WSPatternHandle struct {
	Handler    reflect.Type
	requestUrl string
	config     WSHandlerConfig
}

// SetConfig 设置 WebSocket handler 配置
func (wsHandler *WSPatternHandle) SetConfig(config WSHandlerConfig) {
	wsHandler.config = config
}

// Handle 处理 WebSocket 连接请求
// 功能:
//
//	将 HTTP 连接升级为 WebSocket 连接，加载对应的处理器并执行通信逻辑
//	如果升级失败会记录错误日志并返回错误响应给客户端
//	连接成功后自动调用处理器的 OnConnect 钩子，然后进入通信循环
//
//	- responseWriter: HTTP 响应写入器
//	- request: HTTP 请求对象
func (wsHandler *WSPatternHandle) Handle(responseWriter http.ResponseWriter, request *http.Request) {
	// 创建 upgrader 实例
	upgrader := &websocket.Upgrader{
		CheckOrigin:       wsHandler.config.CheckOrigin,
		HandshakeTimeout:  wsHandler.config.HandshakeTimeout,
		ReadBufferSize:    wsHandler.config.ReadBufferSize,
		WriteBufferSize:   wsHandler.config.WriteBufferSize,
		EnableCompression: wsHandler.config.EnableCompression,
	}

	// 记录客户端信息
	clientIP, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		clientIP = request.RemoteAddr
	}
	userAgent := request.UserAgent()

	// 1. 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		// 记录详细错误日志
		logger.Error.Printf("WebSocket upgrade failed | URL: %s | IP: %s | UserAgent: %s | Error: %v",
			wsHandler.requestUrl, clientIP, userAgent, err)

		// 尝试向客户端返回错误响应
		if responseWriter.Header().Get("Content-Type") == "" {
			responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("WebSocket upgrade failed: " + err.Error()))
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Warning.Printf("WebSocket close error | URL: %s | IP: %s | Error: %v",
				wsHandler.requestUrl, clientIP, err)
		} else {
			logger.Info.Printf("WebSocket disconnected | URL: %s | IP: %s | UserAgent: %s",
				wsHandler.requestUrl, clientIP, userAgent)
		}
	}()

	logger.Info.Printf("WebSocket connected | URL: %s | IP: %s | UserAgent: %s",
		wsHandler.requestUrl, clientIP, userAgent)

	// 2. 创建 handler 实例
	handler := reflect.New(wsHandler.Handler)

	// 3. 检查 handler 是否实现了必要的接口
	// 如果 handler 有 conn 字段，使用反射设置
	if connField := handler.Elem().FieldByName("conn"); connField.IsValid() && connField.CanSet() {
		connField.Set(reflect.ValueOf(conn))
	}

	// 如果 handler 有 request 字段，设置它
	if reqField := handler.Elem().FieldByName("request"); reqField.IsValid() && reqField.CanSet() {
		reqField.Set(reflect.ValueOf(request))
	}

	// 如果 handler 有 SetRequest 方法，调用它
	if setReqMethod := handler.MethodByName("SetRequest"); setReqMethod.IsValid() {
		setReqMethod.Call([]reflect.Value{reflect.ValueOf(request)})
	}

	// 如果 handler 有 SetConn 方法，调用它
	if setConnMethod := handler.MethodByName("SetConn"); setConnMethod.IsValid() {
		setConnMethod.Call([]reflect.Value{reflect.ValueOf(conn)})
	}

	// 如果 handler 有 SetConfig 方法，调用它设置 WSConfig
	if setConfigMethod := handler.MethodByName("SetConfig"); setConfigMethod.IsValid() {
		wsConfig := WSConfig{
			HandshakeTimeout:  wsHandler.config.HandshakeTimeout,
			ReadBufferSize:    wsHandler.config.ReadBufferSize,
			WriteBufferSize:   wsHandler.config.WriteBufferSize,
			EnableCompression: wsHandler.config.EnableCompression,
			CheckOrigin:       wsHandler.config.CheckOrigin,
		}
		setConfigMethod.Call([]reflect.Value{reflect.ValueOf(wsConfig)})
	}

	// 4. 调用 Communicate 方法开始通信
	VoidFuncCall(handler, FnWebSocketCommunicate)
}
