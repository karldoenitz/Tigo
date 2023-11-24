package TigoWeb

// SessionInterface Tigo的session接口，第三方需要自定session底层的实现，每个SessionInterface必须包含NewSessionManager()函数。
// 不想自己实现的可以看看这个插件，直接引入就能用。https://github.com/karldoenitz/tission
type SessionInterface interface {
	// NewSessionManager 新建一个SessionManager
	NewSessionManager() SessionManager
}

// SessionManager session管理器，对session进行生成、获取、删除操作。
type SessionManager interface {
	// GenerateSession 生成session
	GenerateSession(expire int) Session
	// GetSessionBySid 根据session id获取session
	GetSessionBySid(sid string) Session
	// DeleteSession 根据session id删除session
	DeleteSession(sid string)
}

// Session session接口，通过该类型的实例进行session的增删改查。
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

var GlobalSessionManager SessionManager // 全局session管理器
var SessionCookieName = "TigoSessionId" // session的cookie名称
