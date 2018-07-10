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

```
### func (*BaseHandler)GetParameter
```go

```
### func (*BaseHandler)GetHeader
```go

```
### func (*BaseHandler)SetHeader
```go

```
### func (*BaseHandler)GetCookie
```go

```
### func (*BaseHandler)SetCookie
```go

```
### func (*BaseHandler)GetSecureCookie
```go

```
### func (*BaseHandler)SetSecureCookie
```go

```
### func (*BaseHandler)GetCookieObject
```go

```
### func (*BaseHandler)SetCookieObject
```go

```
### func (*BaseHandler)ClearCookie
```go

```
### func (*BaseHandler)ClearAllCookie
```go

```
### func (*BaseHandler)CheckRequestMethodValid
```go

```
### func (*BaseHandler)Redirect
```go

```
### func (*BaseHandler)RedirectPermanently
```go

```
### func (*BaseHandler)ResponseAsHtml
```go

```
### func (*BaseHandler)ResponseAsText
```go

```
### func (*BaseHandler)ResponseAsJson
```go

```
### func (*BaseHandler)ToJson
```go

```

# Tigo.logger
