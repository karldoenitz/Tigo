# 文档([For English Documentation Click Here](https://github.com/karldoenitz/Tigo/blob/master/documentation_en.md))
Tigo是一款用go开发的web应用框架，基于net/http库实现，目的是用来快速搭建restful服务。
# Tigo.TigoWeb
TigoWeb是Tigo框架中的核心部分，Handler、URLpattern以及Application三大核心组件包含于此。
## type BaseHandler
```go
type BaseHandler struct {
    ResponseWriter  http.ResponseWriter
    Request        *http.Request
}
```
```BaseHandler```是一切handler的父结构体，开发者开发的handler必须继承此结构体。
### func (*BaseHandler)InitHandler
```go
func (baseHandler *BaseHandler)InitHandler(responseWriter http.ResponseWriter, request *http.Request)
```
```InitHandler```方法是初始化结构体必须要使用的方法，所有的Handler中的Handle方法必须调用此方法。
### func (*BaseHandler)GetJsonValue
```go
func (baseHandler *BaseHandler)GetJsonValue(key string) (interface{})
```
```GetJsonValue```方法是根据key获取客户端传递的json对象。客户端发送的请求的Content-Type必须是application/json。
### func (*BaseHandler)GetParameter
```go
func (baseHandler *BaseHandler)GetParameter(key string) (value string)
```
```GetParameter```方法是根据key获取客户端传递的参数值。
- 可以从URL上获取参数值，或是form中；如果http的body是json，也可以从中解析出参数值。
### func (*BaseHandler)GetHeader
```go
func (baseHandler *BaseHandler)GetHeader(name string) (value string)
```
```GetHeader```方法是根据name获取http的header值。
### func (*BaseHandler)SetHeader
```go
func (baseHandler *BaseHandler)SetHeader(name string, value string)
```
```SetHeader```方法是根据name设置http的header值。
### func (*BaseHandler)GetCookie
```go
func (baseHandler *BaseHandler)GetCookie(name string) (value string)
```
```GetCookie```方法是根据name获取cookie值。
### func (*BaseHandler)SetCookie
```go
func (baseHandler *BaseHandler)SetCookie(name string, value string)
```
```SetCookie```方法是根据name设置cookie值。
### func (*BaseHandler)GetSecureCookie
```go
func (baseHandler *BaseHandler)GetSecureCookie(name string, key ...string) (value string)
```
```GetSecureCookie```方法是用来获取加密的cookie值，key为解密所需要用到的密钥，key可以不填，也可以在configuration配置文件中配置。
- 当configuration文件和函数参数中都设置了key，则以函数中设置的key为准；
- 当configuration文件和函数参数中都未设置key，则依然会进行加密。
### func (*BaseHandler)SetSecureCookie
```go
func (baseHandler *BaseHandler)SetSecureCookie(name string, value string, key ...string)
```
```SetSecureCookie```方法是用来给客户端设置一个加密cookie，key为加密所需要用到的密钥，key可以不填，也可以在configuration配置文件中配置。
- 当configuration文件和函数参数中都设置了key，则以函数中设置的key为准；
- 当configuration文件和函数参数中都未设置key，则依然会进行加密。
### func (*BaseHandler)GetCookieObject
```go
func (baseHandler *BaseHandler)GetCookieObject(name ...string) (Cookie, error)
```
```GetCookieObject```方法是用来根据name获取指定的cookie对象，返回值类型为```Cookie```。
### func (*BaseHandler)SetCookieObject
```go
func (baseHandler *BaseHandler)SetCookieObject(cookie Cookie)
```
```SetCookieObject```方法接收一个```Cookie```对象，为客户端设置cookie。
### func (*BaseHandler)ClearCookie
```go
func (baseHandler *BaseHandler)ClearCookie(name string)
```
```ClearCookie```方法是根据指定的name清除cookie。
### func (*BaseHandler)ClearAllCookie
```go
func (baseHandler *BaseHandler)ClearAllCookie()
```
```ClearAllCookie```方法是用来清空所有的cookie的。
### func (*BaseHandler)CheckRequestMethodValid
```go
func (baseHandler *BaseHandler)CheckRequestMethodValid(methods ...string) (result bool)
```
```CheckRequestMethodValid```方法是用来判断http请求方法是否合法，接收多个参数，参数表示支持的http请求方式，不填参数则认为不支持所有请求方式。
### func (*BaseHandler)Redirect
```go
func (baseHandler *BaseHandler)Redirect(url string, expire ...time.Time)
```
```Redirect```方法是将当前handler所挂载的URL重定向到另一个URL地址，expire为客户端过期时间，如果不填写expire，则会使用客户端默认过期时间。
### func (*BaseHandler)RedirectPermanently
```go
func (baseHandler *BaseHandler)RedirectPermanently(url string)
```
```RedirectPermanently```方法是将当前handler所挂载的URL永久性重定向到另一个URL地址。
### func (*BaseHandler)ResponseAsHtml
```go
func (baseHandler *BaseHandler)ResponseAsHtml(result string)
```
```ResponseAsHtml```方法是将传入的字符串以html文本类型响应给客户端。
### func (*BaseHandler)ResponseAsText
```go
func (baseHandler *BaseHandler)ResponseAsText(result string)
```
```ResponseAsText```方法是将传入的字符串以text文本类型响应给客户端。
### func (*BaseHandler)ResponseAsJson
```go
func (baseHandler *BaseHandler)ResponseAsJson(response Response)
```
```ResponseAsJson```方法是将传入的对象转换成json字符串，然后响应给客户端，如果转换失败则会向客户端响应空字符串。
### func (*BaseHandler)ToJson
```go
func (baseHandler *BaseHandler)ToJson(response Response) (result string)
```
```ToJson```方法是将一个对象转换成json字符串。如果转换失败则会返回空字符串。
## type UrlPattern
```go
type UrlPattern struct {
    UrlMapping map[string] interface{Handle(http.ResponseWriter, *http.Request)}
}
```
URL路由设置，使用这个结构体在应用中配置URL与对应的handler。
### func (urlPattern *UrlPattern)AppendUrlPattern
```go
func (urlPattern *UrlPattern)AppendUrlPattern(uri string, v interface{Handle(http.ResponseWriter, *http.Request)})
```
此方法是向指定URL上挂载一个Handler。
### func (urlPattern *UrlPattern)Init
```go
func (urlPattern *UrlPattern)Init()
```
初始化URL映射。
# Tigo.logger
