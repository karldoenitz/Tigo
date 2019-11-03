package TigoWeb

type SessionInterface interface {
	// NewSessionManager 新建一个SessionManager
	NewSessionManager() SessionManager
}

type SessionManager interface {
	// GenerateSession 生成session
	GenerateSession(expire int) Session
	// GetSessionBySid 根据session id获取session
	GetSessionBySid(sid string) Session
	// DeleteSession 根据session id删除session
	DeleteSession(sid string)
}

type Session interface {
	// Set 设置session值，设置失败返回error
	Set(key string, value interface{}) error
	// Get 获取session值，获取失败返回error
	Get(key string, value interface{}) error
	// Delete 删除session
	Delete(key string)
	// SessionId 获取session id
	SessionId() string
}

var GlobalSessionManager SessionManager
var SessionCookieName = "TigoSessionId"
