package TigoWeb

type SessionInterface interface {
	// 新建一个SessionManager
	NewSessionManager() SessionManager
}

type SessionManager interface {
	// 生成session
	GenerateSession(expire int) Session
	// 根据session id获取session
	GetSessionBySid(sid string) Session
	// 根据session id删除session
	DeleteSession(sid string)
}

type Session interface {
	// 设置session值，设置失败返回error
	Set(key string, value interface{}) error
	// 获取session值，获取失败返回error
	Get(key string, value interface{}) error
	// 删除session
	Delete(key string)
	// 获取session id
	SessionId() string
}

var GlobalSessionManager SessionManager
var SessionCookieName = "TigoSessionId"
