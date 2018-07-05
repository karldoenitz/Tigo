package WebFramework

import "time"

// 自定义Cookie结构体，可参看http.Cookie
type Cookie struct {
	Name        string
	Value       string

	IsSecurity  bool      // 是否对cookie值进行加密
	SecurityKey string    // 加密cookie用到的key

	Path        string    // 可选
	Domain      string    // 可选
	Expires     time.Time // 可选
	RawExpires  string    // 只有在读取Cookie时有效

	// MaxAge=0 表示未指定“Max-Age”属性
	// MaxAge<0 表示现在删除cookie，相当于'Max-Age：0'
	// MaxAge>0 表示Max-Age属性存在并以秒为单位给出
	MaxAge      int
	Secure      bool
	HttpOnly    bool
	Raw         string
	Unparsed  []string    // 原始文本中未解析的属性值
}

// 获取cookie加密值
//   - IsSecurity如果设置为false，则返回原始值
//   - IsSecurity如果设置为true，则返回加密后的值
// 如果加密失败，则抛出异常
func (cookie *Cookie)GetCookieEncodeValue()(result string) {
	if !cookie.IsSecurity {
		return cookie.Value
	}
	value := []byte(cookie.Value)
	key   := []byte(cookie.SecurityKey)
	securityValue, err := DesEncrypt(value, key)
	if err != nil {
		panic(err)
	}
	result = string(securityValue)
	return result
}
