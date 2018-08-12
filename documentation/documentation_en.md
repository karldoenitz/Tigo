# Documentation([中文文档点击此处](https://github.com/karldoenitz/Tigo/blob/master/documentation/documentation.md))
Tigo is a web framework developed in go language, based on net/http. In order to build web server quickly.  
API index:
- [TigoWeb](#TigoWeb)
  - [type BaseHandler](#BaseHandler)
    - [func InitHandler](#InitHandler)
    - [func GetJsonValue](#GetJsonValue)
    - [func GetParameter](#GetParameter)
    - [func GetHeader](#GetHeader)
    - [func SetHeader](#SetHeader)
    - [func GetCookie](#GetCookie)
    - [func SetCookie](#SetCookie)
    - [func GetSecureCookie](#GetSecureCookie)
    - [func SetSecureCookie](#SetSecureCookie)
    - [func GetCookieObject](#GetCookieObject)
    - [func SetCookieObject](#SetCookieObject)
    - [func ClearCookie](#ClearCookie)
    - [func ClearAllCookie](#ClearAllCookie)
    - [func CheckRequestMethodValid](#CheckRequestMethodValid)
    - [func Redirect](#Redirect)
    - [func RedirectPermanently](#RedirectPermanently)
    - [func ResponseAsHtml](#ResponseAsHtml)
    - [func ResponseAsText](#ResponseAsText)
    - [func ResponseAsJson](#ResponseAsJson)
    - [func ToJson](#ToJson)
  - [type UrlPattern](#UrlPattern)
    - [func AppendUrlPattern](#AppendUrlPattern)
    - [func Init](#Init)
  - [type Application](#Application)
    - [func Run](#Run)
  - [type Cookie](#Cookie)
    - [func GetCookieEncodeValue](#GetCookieEncodeValue)
    - [func GetCookieDecodeValue](#GetCookieDecodeValue)
    - [func ToHttpCookie](#ToHttpCookie)
    - [func ConvertFromHttpCookie](#ConvertFromHttpCookie)
    - [func SetSecurityKey](#SetSecurityKey)
  - [type BaseResponse](#BaseResponse)
    - [func Print](#Print)
    - [func ToJson](#ResponseToJson)
  - [type GlobalConfig](#GlobalConfig)
    - [func Init](#GlobalInit)
  - [utils](#utils)
    - [func Encrypt](#Encrypt)
    - [func Decrypt](#Decrypt)
    - [func InitGlobalConfig](#InitGlobalConfig)
- [logger](#logger)
  - [Demo](#logDemo)
  - [structure](#LogStructure)
    - [type LogLevel](#LogLevel)
  - [functions](#loggerFunctions)
    - [func SetLogPath](#SetLogPath)
    - [func InitLoggerWithConfigFile](#InitLoggerWithConfigFile)
    - [func InitLoggerWithObject](#InitLoggerWithObject)
    - [func InitTrace](#InitTrace)
    - [func InitInfo](#InitInfo)
    - [func InitWarning](#InitWarning)
    - [func InitError](#InitError)
# Tigo.TigoWeb<a name="TigoWeb"></a>
TigoWeb is the core part of Tigo framework, it contains Handler,URLpattern and Application.
## type BaseHandler<a name="BaseHandler"></a>
```go
type BaseHandler struct {
    ResponseWriter  http.ResponseWriter
    Request        *http.Request
}
```
```BaseHandler``` is the base structure of all handlers, the handlers developed by developers must extend this structure.
### func (*BaseHandler)InitHandler<a name="InitHandler"></a>
```go
func (baseHandler *BaseHandler)InitHandler(responseWriter http.ResponseWriter, request *http.Request)
```
```InitHandler``` is the method to initialize handler, all the ```Handle``` method must call this method in handlers.
### func (*BaseHandler)GetJsonValue<a name="GetJsonValue"></a>
```go
func (baseHandler *BaseHandler)GetJsonValue(key string) (interface{})
```
```GetJsonValue``` is the method to get json object value in http body, the request's Content-Type must be application/json.
### func (*BaseHandler)GetParameter<a name="GetParameter"></a>
```go
func (baseHandler *BaseHandler)GetParameter(key string) (value string)
```
```GetParameter``` is the method to get value by key.
- Get the value in url, form or json.
### func (*BaseHandler)GetHeader<a name="GetHeader"></a>
```go
func (baseHandler *BaseHandler)GetHeader(name string) (value string)
```
```GetHeader``` is the method to get http header.
### func (*BaseHandler)SetHeader<a name="SetHeader"></a>
```go
func (baseHandler *BaseHandler)SetHeader(name string, value string)
```
```SetHeader``` is the method to set http header.
### func (*BaseHandler)GetCookie<a name="GetCookie"></a>
```go
func (baseHandler *BaseHandler)GetCookie(name string) (value string)
```
```GetCookie``` is the method to get cookie value.
### func (*BaseHandler)SetCookie<a name="SetCookie"></a>
```go
func (baseHandler *BaseHandler)SetCookie(name string, value string)
```
```SetCookie``` is the method to set cookie value.
### func (*BaseHandler)GetSecureCookie<a name="GetSecureCookie"></a>
```go
func (baseHandler *BaseHandler)GetSecureCookie(name string, key ...string) (value string)
```
```GetSecureCookie``` is the method to get cookie's security value.
- The key pass to method will instead of the key in configuration;
- This method can decrypt cookie value without key.
### func (*BaseHandler)SetSecureCookie<a name="SetSecureCookie"></a>
```go
func (baseHandler *BaseHandler)SetSecureCookie(name string, value string, key ...string)
```
```SetSecureCookie``` is the method to set cookie's security value.
- The key pass to method will instead of the key in configuration;
- This method can encrypt cookie value without key.
### func (*BaseHandler)GetCookieObject<a name="GetCookieObject"></a>
```go
func (baseHandler *BaseHandler)GetCookieObject(name ...string) (Cookie, error)
```
```GetCookieObject``` is the method to get cookie object by name.
### func (*BaseHandler)SetCookieObject<a name="SetCookieObject"></a>
```go
func (baseHandler *BaseHandler)SetCookieObject(cookie Cookie)
```
```SetCookieObject``` is the method to set cookie object by name.
### func (*BaseHandler)ClearCookie<a name="ClearCookie"></a>
```go
func (baseHandler *BaseHandler)ClearCookie(name string)
```
```ClearCookie``` is the method to clear cookie.
### func (*BaseHandler)ClearAllCookie<a name="ClearAllCookie"></a>
```go
func (baseHandler *BaseHandler)ClearAllCookie()
```
```ClearAllCookie``` is the method to clear all cookie.
### func (*BaseHandler)CheckRequestMethodValid<a name="CheckRequestMethodValid"></a>
```go
func (baseHandler *BaseHandler)CheckRequestMethodValid(methods ...string) (result bool)
```
```CheckRequestMethodValid``` is the method to check request method, example:
```go
func (h *Handler)Handle(responseWriter http.ResponseWriter, request *http.Request) {
    h.InitHandler(responseWriter, request)
    if !h.CheckRequestMethodValid("GET", "POST") {
        return
    }
}
```
### func (*BaseHandler)Redirect<a name="Redirect"></a>
```go
func (baseHandler *BaseHandler)Redirect(url string, expire ...time.Time)
```
```Redirect``` is the method to redirect client to another url.
### func (*BaseHandler)RedirectPermanently<a name="RedirectPermanently"></a>
```go
func (baseHandler *BaseHandler)RedirectPermanently(url string)
```
```RedirectPermanently``` is the method to redirect client to another url permanently.
### func (baseHandler *BaseHandler)Render<a name="Render"></a>
```go
func (baseHandler *BaseHandler)Render(data interface{}, templates ...string)
```
`Render` is the method to compile html template and response the result to client.

parameters:
- data: the instance of any type structure;
- templates: the name of html templates.

**Attention:**  
If not configure base path of templates in configuration, this method will use relative path.
### func (*BaseHandler)ResponseAsHtml<a name="ResponseAsHtml"></a>
```go
func (baseHandler *BaseHandler)ResponseAsHtml(result string)
```
```ResponseAsHtml``` is the method to response with html.
### func (*BaseHandler)ResponseAsText<a name="ResponseAsText"></a>
```go
func (baseHandler *BaseHandler)ResponseAsText(result string)
```
```ResponseAsText``` is the method to response with text.
### func (*BaseHandler)ResponseAsJson<a name="ResponseAsJson"></a>
```go
func (baseHandler *BaseHandler)ResponseAsJson(response Response)
```
```ResponseAsJson``` is the method to response with json, this method will response empty string if convert response object to json failed.
### func (*BaseHandler)ToJson<a name="ToJson"></a>
```go
func (baseHandler *BaseHandler)ToJson(response Response) (result string)
```
```ToJson``` is the method to convert response object to json string.
## type UrlPattern<a name="UrlPattern"></a>
```go
type UrlPattern struct {
    UrlMapping map[string] interface{Handle(http.ResponseWriter, *http.Request)}
}
```
URL mapping structure.
### func (urlPattern *UrlPattern)AppendUrlPattern<a name="AppendUrlPattern"></a>
```go
func (urlPattern *UrlPattern)AppendUrlPattern(uri string, v interface{Handle(http.ResponseWriter, *http.Request)})
```
Use this method to append a handler to an url.
### func (urlPattern *UrlPattern)Init<a name="Init"></a>
```go
func (urlPattern *UrlPattern)Init()
```
Initialize url mapping.
## type Application<a name="Application"></a>
```go
type Application struct {
    IPAddress  string      // IP address
    Port       int         // port
    UrlPattern UrlPattern  // url mapping
    ConfigPath string      // global config
}
```
Application is the launcher of web server developed by Tigo.
- IPAddress: IP address, developers can config IP address in configuration file;
- Port: Port, developers can configure it in configuration file;
- UrlPattern: url mapping.
- ConfigPath: the configuration's path.
### func (application *Application)Run<a name="Run"></a>
```go
func (application *Application)Run()
```
Run the web server use this method, if config `cert` and `key` file, the server will run with https protocol.
## type Cookie<a name="Cookie"></a>
```go
type Cookie struct {
    Name        string
    Value       string

    IsSecurity  bool      // whether encrypt cookie's value
    SecurityKey string    // the key encrypt cookie's value

    Path        string    // options
    Domain      string    // options
    Expires     time.Time // options
    RawExpires  string    // for reading cookies only

    // MaxAge=0 means no 'Max-Age' attribute specified.
    // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
    // MaxAge>0 means Max-Age attribute present and given in seconds
    MaxAge      int
    Secure      bool
    HttpOnly    bool
    Raw         string
    Unparsed  []string    // Raw text of unparsed attribute-value pairs
}
```
The cookie structure, use this structure to operate the cookie.
### func (cookie *Cookie)GetCookieEncodeValue<a name="GetCookieEncodeValue"></a>
```go
func (cookie *Cookie)GetCookieEncodeValue()(result string)
```
Use this method to get cookie's encode value.
### func (cookie *Cookie)GetCookieDecodeValue<a name="GetCookieDecodeValue"></a>
```go
func (cookie *Cookie)GetCookieDecodeValue()(result string)
```
Use this method to get cookie's decode value.
### func (cookie *Cookie)ToHttpCookie<a name="ToHttpCookie"></a>
```go
func (cookie *Cookie)ToHttpCookie()(http.Cookie)
```
Convert cookie to http.Cookie.
### func (cookie *Cookie)ConvertFromHttpCookie<a name="ConvertFromHttpCookie"></a>
```go
func (cookie *Cookie)ConvertFromHttpCookie(httpCookie http.Cookie)
```
Convert http.Cookie to cookie.
### func (cookie *Cookie)SetSecurityKey<a name="SetSecurityKey"></a>
```go
func (cookie *Cookie)SetSecurityKey(key string)
```
Set security cookie's key.
## type BaseResponse<a name="BaseResponse"></a>
```go
type BaseResponse struct {

}
```
The structure extend this structure can use method ```func (baseHandler *BaseHandler)ResponseAsJson(response Response)``` to response json to client.
### func (baseResponse *BaseResponse)Print<a name="Print"></a>
```
func (baseResponse *BaseResponse)Print()
```
Use this method to print json to console.
### func (baseResponse *BaseResponse)ToJson<a name="ResponseToJson"></a>
```
func (baseResponse *BaseResponse)ToJson() (string)
```
Use this method to convert object to json string.
## type GlobalConfig<a name="GlobalConfig"></a>
```go
type GlobalConfig struct {
    IP       string           `json:"ip"`        // IP address
    Port     int              `json:"port"`      // Port
    Cert     string           `json:"cert"`      // https cert
    CertKey  string           `json:"cert_key"`  // https key
    Cookie   string           `json:"cookie"`    // cookie's security key
    Log      logger.LogLevel  `json:"log_path"`  // log config
}
```
### func (globalConfig *GlobalConfig)Init<a name="GlobalInit"></a>
```go
func (globalConfig *GlobalConfig)Init(configPath string)
```
Use configuration file to initialize global configuration.    
Glance:
- IP: IP address
- Port: the Port server blinded
- Cert: the path of https cert file
- CertKey: the path of https key file
- Cookie: the key to encrypt/decrypt cookie
- Log: the instance of Tigo.logger.LogLevel

configuration.json example：
```javascript
{
    "ip": "0.0.0.0",
    "port": 8888,
    "cert": "/home/work/certfile.ext"
    "cert_key": "/home/work/certkeyfile.ext",
    "log_path": "/home/work/log/server_run.log",
    "cookie": "thisiscookiekey"
}
```
## utils<a name="utils"></a>
Encrypt method<a name="Encrypt"></a>
```go
func Encrypt(src[]byte, key []byte) string
```
Use this method to encrypt byte array.  
Decrypt method<a name="Decrypt"></a>
```go
func Decrypt(src[]byte, key []byte) ([]byte)
```
Use this method to decrypt byte array.  
Initialize global configuration<a name="InitGlobalConfig"></a>
```go
func InitGlobalConfig(configPath string)
```
Use this method to initialize global configuration.
# Tigo.logger<a name="logger"></a>
Use this package to print log.  
## Demo<a name="logDemo"></a>
The demo about using ```logger``` package in the web application developed by Tigo.
```go
package main

import (
    "net/http"
    "github.com/karldoenitz/Tigo/TigoWeb"
    "github.com/karldoenitz/Tigo/logger"
)

type HelloHandler struct {
    TigoWeb.BaseHandler
}

func (helloHandler *HelloHandler)Handle() {
    logger.Info.Printf("info data: %s", "test") // print log here
    helloHandler.ResponseAsHtml("<p1 style='color: red'>Hello Tigo!</p1>")
}

var urls = map[string]interface{}{
    "/hello-tigo": &HelloHandler{},
}

func main() {
    application := TigoWeb.Application{
        UrlPattern: urls,
        ConfigPath: "./configuration.json",  // use configuration file to config logger, or you can use LogLevel instance if you wanted.
    }
    application.Run()
}
```
```configuration.json``` content：
```JavaScript
{
  "cookie": "TencentCode",
  "ip": "0.0.0.0",
  "port": 8080,
  "log": {
    "trace": "stdout",  // only print trace message in console
    "info": "/Users/karllee/Desktop/run-info.log",  // print info message in run-info.log file
    "warning": "/Users/karllee/Desktop/run.log",  // print warning info and error info in same file
    "error": "/Users/karllee/Desktop/run.log"
  }
}
```
If you wanna use logger package in your application not base on Tigo, you can use `func (globalConfig *GlobalConfig)Init(configPath string)` function to initialize logger package.
## Structure<a name="LogStructure"></a>
The structure in logger package:
### type LogLevel<a name="LogLevel"></a>
```go
// log level structure
//   - Trace
//   - Info
//   - Warning
//   - Error
// discard: discard, stdout: print in console; the path of log file
type LogLevel struct {
    Trace    string   `json:"trace"`
    Info     string   `json:"info"`
    Warning  string   `json:"warning"`
    Error    string   `json:"error"`
}
```
Initialize this structure and pass the instance to ```InitLoggerWithObject``` to init logger package.
## logger inner functions<a name="loggerFunctions"></a>
### func SetLogPath<a name="SetLogPath"></a>
Set Log file's Path<a name="SetLogPath"></a>
```go
func SetLogPath(logPath string)
```
example：
```go
import "github.com/karldoenitz/Tigo/logger"

logger.Info.Printf("It is a test...")
logger.Warning.Printf("warning!")
logger.Error.Printf("ERROR!!!")
```
Attention: this method can override the config you set before.
### func InitLoggerWithConfigFile<a name="InitLoggerWithConfigFile"></a>
```go
func InitLoggerWithConfigFile(filePath string)
```
Initialize logger with configuration file.
### func InitLoggerWithObject<a name="InitLoggerWithObject"></a>
```go
func InitLoggerWithObject(logLevel LogLevel)
```
Initialize logger with LogLevel instance.
### func InitTrace<a name="InitTrace"></a>
```go
func InitTrace(level string)
```
Initialize Trace instance.  
Parameter values:
- discard: discard log message;
- stdout: print log message in console;
- real file path: the path of log file.
### func InitInfo<a name="InitInfo"></a>
```go
func InitInfo(level string)
```
Initialize Info instance.  
Parameter values:
- discard: discard log message;
- stdout: print log message in console;
- real file path: the path of log file.
### func InitWarning<a name="InitWarning"></a>
```go
func InitWarning(level string)
```
Initialize Warning instance.  
Parameter values:
- discard: discard log message;
- stdout: print log message in console;
- real file path: the path of log file.
### func InitError<a name="InitError"></a>
```go
func InitError(level string)
```
Initialize Error instance.  
Parameter values:
- discard: discard log message;
- stdout: print log message in console;
- real file path: the path of log file.
