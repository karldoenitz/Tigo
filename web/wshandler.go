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

type WSPatternHandle struct {
	Handler    reflect.Type
	requestUrl string
}

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
	logger.Info.Printf("websocket | %s | connect success", wsHandler.requestUrl)
	// 调用 Communicate 方法
	VoidFuncCall(handler, FnWebSocketCommunicate, reflect.ValueOf(conn))
}

type WSBaseHandler struct {
}

func (h *WSBaseHandler) Communicate(conn *websocket.Conn) {
	respMsg := []byte("Communicate is not implemented")
	if err := conn.WriteMessage(websocket.TextMessage, respMsg); err != nil {
		logger.Error.Println("发送错误:", err)
	}
}
