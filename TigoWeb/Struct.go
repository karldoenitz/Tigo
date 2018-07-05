package WebFramework

import (
	"time"
	"net/http"
)

//////////////////////////////////////Structure Cookie//////////////////////////////////////////////////////////////////

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

// 获取cookie解密值
//   - IsSecurity如果设置为false，则返回原始值
//   - IsSecurity如果设置为true，则返回加密后的值
// 如果解密失败，则抛出异常
func (cookie *Cookie)GetCookieDecodeValue()(result string) {
	if !cookie.IsSecurity {
		return cookie.Value
	}
	value := []byte(cookie.Value)
	key   := []byte(cookie.SecurityKey)
	securityValue, err := DesDecrypt(value, key)
	if err != nil {
		panic(err)
	}
	result = string(securityValue)
	return result
}

// 转换为http/Cookie对象
func (cookie *Cookie)ToHttpCookie()(http.Cookie) {
	httpCookie := http.Cookie{
		Name:       cookie.Name,
		Value:      cookie.GetCookieEncodeValue(),
		Path:       cookie.Path,
		Domain:     cookie.Domain,
		Expires:    cookie.Expires,
		RawExpires: cookie.RawExpires,
		MaxAge:     cookie.MaxAge,
		Secure:     cookie.Secure,
		HttpOnly:   cookie.HttpOnly,
		Raw:        cookie.Raw,
		Unparsed:   cookie.Unparsed,
	}
	return httpCookie
}

// 将http/Cookie转换为Cookie
func (cookie *Cookie)ConvertFromHttpCookie(httpCookie http.Cookie) {
	cookie.Name         = httpCookie.Name
	cookie.Value        = httpCookie.Value

	cookie.Path         = httpCookie.Path
	cookie.Domain       = httpCookie.Domain
	cookie.Expires      = httpCookie.Expires
	cookie.RawExpires   = httpCookie.RawExpires

	cookie.MaxAge       = httpCookie.MaxAge
	cookie.Secure       = httpCookie.Secure
	cookie.HttpOnly     = httpCookie.HttpOnly
	cookie.Raw          = httpCookie.Raw
	cookie.Unparsed     = httpCookie.Unparsed
}

// 为Cookie设置SecurityKey
func (cookie *Cookie)SetSecurityKey(key string) {
	cookie.SecurityKey = key
	cookie.IsSecurity = true
}
