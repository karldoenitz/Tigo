# Documentation([中文文档点击此处](https://github.com/karldoenitz/Tigo/blob/master/documentation/documentation.md))
Tigo is a web framework developed in go language, based on net/http. In order to build web server quickly.  
API index:
- [TigoWeb](#tigotigoweb)
  - [type BaseHandler](#type-basehandler)
    - [func InitHandler](#func-basehandlerinithandler)
    - [func GetJsonValue](#func-basehandlergetjsonvalue)
    - [func GetParameter](#func-basehandlergetparameter)
    - [func GetHeader](#func-basehandlergetheader)
    - [func SetHeader](#func-basehandlersetheader)
    - [func GetCtxVal](#func-basehandlersetctxval)
    - [func SetCtxVal](#func-basehandlersetctxval)
    - [func GetCookie](#func-basehandlergetcookie)
    - [func SetCookie](#func-basehandlersetcookie)
    - [func GetSecureCookie](#func-basehandlergetsecurecookie)
    - [func SetSecureCookie](#func-basehandlersetsecurecookie)
    - [func GetCookieObject](#func-basehandlergetcookieobject)
    - [func SetCookieObject](#func-basehandlersetcookieobject)
    - [func SetAdvancedCookie](#func-basehandlersetadvancedcookie)
    - [func ClearCookie](#func-basehandlerclearcookie)
    - [func ClearAllCookie](#func-basehandlerclearallcookie)
    - [func Redirect](#func-basehandlerredirect)
    - [func RedirectPermanently](#func-basehandlerredirectpermanently)
    - [func Render](#func-basehandler-basehandlerrender)
    - [func ResponseAsHtml](#func-basehandlerresponseashtml)
    - [func ResponseAsText](#func-basehandlerresponseastext)
    - [func ResponseAsJson](#func-basehandlerresponseasjson)
    - [func ResponseFmt](#func-basehandler-responsefmt)
    - [func ServerError](#func-basehandler-servererror)
    - [func ToJson](#func-basehandlertojson)
    - [func DumpHttpRequestMsg](#func-basehandlerdumphttprequestmsg)
    - [func CheckJsonBinding](#func-basehandlercheckjsonbinding)
    - [func CheckFormBinding](#func-basehandlercheckformbinding)
    - [func CheckUrlParamBinding](#func-basehandlercheckurlparambinding)
    - [func CheckParamBinding](#func-basehandlercheckparambinding)
    - [func BeforeRequest](#func-basehandlerbeforerequest)
    - [func TeardownRequest](#func-basehandlerteardownrequest)
  - [type UrlPattern](#type-urlpattern)
    - [func AppendUrlPattern](#func-urlpattern-urlpatternappendurlpattern)
    - [func AppendRouterPattern](#func-urlpattern-urlpattern-appendrouterpattern)
    - [func Init](#func-basehandlerinithandler)
  - [type Application](#type-application)
    - [func Listen](#func-application-applicationlisten)
    - [func MountFileServer](#func-application-applicationmountfileserver)
    - [func Run](#func-application-applicationrun)
  - [type Cookie](#func-basehandlergetcookie)
    - [func GetCookieEncodeValue](#func-cookie-cookiegetcookieencodevalue)
    - [func GetCookieDecodeValue](#func-cookie-cookiegetcookiedecodevalue)
    - [func ToHttpCookie](#func-cookie-cookietohttpcookie)
    - [func ConvertFromHttpCookie](#func-cookie-cookieconvertfromhttpcookie)
    - [func SetSecurityKey](#func-cookie-cookiesetsecuritykey)
  - [type BaseResponse](#type-baseresponse)
    - [func Print](#func-baseresponse-baseresponseprint)
    - [func ToJson](#func-baseresponse-baseresponsetojson)
  - [type GlobalConfig](#type-globalconfig)
    - [func Init](#func-globalconfig-globalconfiginit)
  - [utils](#utils)
    - [func Encrypt](#Encrypt)
    - [func Decrypt](#Decrypt)
    - [func UrlEncode](#urlencode)
    - [func UrlDecode](#urldecode)
    - [func InitGlobalConfig](#InitGlobalConfig)
    - [func InitGlobalConfigWithObj](#InitGlobalConfigWithObj)
- [logger](#tigologger)
  - [Demo](#demo)
  - [structure](#structure)
    - [type LogLevel](#type-loglevel)
  - [functions](#logger-inner-functions)
    - [func SetLogPath](#func-setlogpath)
    - [func InitLoggerWithConfigFile](#func-initloggerwithconfigfile)
    - [func InitLoggerWithObject](#func-initloggerwithobject)
    - [func InitTrace](#func-inittrace)
    - [func InitInfo](#func-initinfo)
    - [func InitWarning](#func-initwarning)
    - [func InitError](#func-initerror)
- [request](#tigorequest)
  - [type Response](#type-response)
    - [ToContentStr](#func-response-responsetocontentstr)
  - [functions](#request-module)
    - [func Request](#func-request)
    - [func MakeRequest](#func-makerequest)
    - [func Get](#func-get)
    - [func Post](#func-post)
    - [func Put](#func-put)
    - [func Patch](#func-patch)
    - [func Head](#func-head)
    - [func Options](#func-options)
    - [func Delete](#func-delete)
- [binding](#tigobinding)
  - [functions](#binding-functions)
    - [func ParseJsonToInstance](#func-parsejsontoinstance)
    - [func ValidateInstance](#func-validateinstance)
    - [func FormBytesToStructure](#func-formbytestostructure)
    - [func ParseFormToInstance](#func-parseformtoinstance)
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
### func (*BaseHandler)GetCtxVal<a name="GetCtxVal"></a>
```go
func (baseHandler *BaseHandler)GetCtxVal(key string) interface{}
```
```GetCtxVal``` get value from http context.
### func (*BaseHandler)SetCtxVal<a name="SetCtxVal"></a>
```go
func (baseHandler *BaseHandler) SetCtxVal(key string, val interface{})
```
```SetCtxVal``` set value to http context.
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
### func (*BaseHandler)SetAdvancedCookie
```go
func (baseHandler *BaseHandler) SetAdvancedCookie(name string, value string, attrs ...string)
```
```SetAdvancedCookie```is the method to set cookie's attributes.  
parameters:
- name cookie's name
- value cookie's value
- attrs cookie's attributes:
  - "path={{string}}" set cookie's path
  - "domain={{string}}" set cookie's domain
  - "raw={{string}}" set cookie's raw
  - "maxAge={{int}}" MaxAge=0 means no 'Max-Age' attribute specified. MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0', MaxAge>0 means Max-Age attribute present and given in seconds
  - "expires={{int}}" set cookie's expires time, given in seconds
  - "secure={{bool}}" set secure
  - "httpOnly={{bool}}" set http only
  - "isSecurity={{bool}}" `true` means cookie will be encrypted 

**Example:**
```go
type DemoHandler struct {
	TigoWeb.BaseHandler
}

func (d *DemoHandler) Get() {
	d.SetAdvancedCookie(
		"test", 
		"value", 
		"path=/", 
		"domain=localhost", 
		"expires=10", 
		"httpOnly=true", 
		"secure=true",
	)
	d.ResponseAsText("test")
}
```
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
func (baseHandler *BaseHandler)ResponseAsHtml(result string, charset ...string)
```
```ResponseAsHtml``` is the method to response with html, charset default `utf-8`.
### func (*BaseHandler)ResponseAsText<a name="ResponseAsText"></a>
```go
func (baseHandler *BaseHandler)ResponseAsText(result string)
```
```ResponseAsText``` is the method to response with text.
### func (*BaseHandler)ResponseAsJson<a name="ResponseAsJson"></a>
```go
func (baseHandler *BaseHandler)ResponseAsJson(response interface{}, charset ...string)
```
```ResponseAsJson``` is the method to response with json, this method will response empty string if convert response object to json failed, charset default `utf-8`.
### func (*BaseHandler) ResponseFmt
```go
func (baseHandler *BaseHandler) ResponseFmt(format string, values... interface{})
```
```ResponseFmt``` is the method to format the result and response to client.
### func (*BaseHandler) ServerError
```go
func (baseHandler *BaseHandler) ServerError(err error)
```
```ServerError``` is the method to response 500 error to client.
### func (*BaseHandler)ToJson<a name="ToJson"></a>
```go
func (baseHandler *BaseHandler)ToJson(response interface{}) (result string)
```
```ToJson``` is the method to convert response object to json string.
### func (*BaseHandler)DumpHttpRequestMsg<a name="DumpHttpRequestMsg"></a>
```go
func (baseHandler *BaseHandler)DumpHttpRequestMsg(logLevel int) (result string)
```
```DumpHttpRequestMsg``` is the method to dump http request message to console or log file.  
The parameter `logLevel`'s value:
- 1: dump message to trace level    // logger.TraceLevel
- 2: dump message to info level     // logger.InfoLevel
- 3: dump message to warning level  // logger.WarningLevel
- 4: dump message to error level    // logger.ErrorLevel
- others: dump message to console
### func (*BaseHandler)CheckJsonBinding<a name="CheckJsonBinding"></a>
```go
func (baseHandler *BaseHandler) CheckJsonBinding(obj interface{}) error
```
```CheckJsonBinding``` check the json from http request.  
tag:
- required: check this field if tag's value is `true`
- default: set default value on this field only `required` is true
- regex: use regular express to check the value of this field
example:
```go
type TestParamCheckHandler struct {
    TigoWeb.BaseHandler
}

func (t *TestParamCheckHandler)Post() {
    params := struct{
        Username string `json:"username" required:"true" regex:"^[0-9a-zA-Z_]{1,}$"`
        Password string `json:"password" required:"true"`
        Age      int    `json:"age" required:"true" default:"18"`
    }{}
    if err := t.CheckJsonBinding(&params); err != nil {
        t.ResponseAsJson(struct{
            Msg string
        }{err.Error()})
        return
    }
    // your program
}
```
Post的json：
```javascript
{
    "username": "wo&ni",
    "password": "tihs si wrodpssa"
} // error "username regex can not match"
{
    "username": "wo_ni",
} // error "password is required"
// age default value is 18
```
reference `Tigo.binding.ValidateInstance`
### func (*BaseHandler)CheckFormBinding<a name="CheckFormBinding"></a>
```go
func (baseHandler *BaseHandler) CheckFormBinding(obj interface{}) error
```
```CheckFormBinding``` check the form from http request.
### func (*BaseHandler)CheckUrlParamBinding
```go
func (baseHandler *BaseHandler) CheckUrlParamBinding(obj interface{}) error
```
```CheckUrlParamBinding``` check the parameters in url.
### func (*BaseHandler)CheckParamBinding<a name="CheckParamBinding"></a>
```go
func (baseHandler *BaseHandler) CheckParamBinding(obj interface{}) error
```
```CheckParamBinding``` check the param from http request, form or json.
### func (*BaseHandler)UrlEncode
```go
func (baseHandler *BaseHandler) UrlEncode(value string) string
```
```UrlEncode``` url encode.
### func (*BaseHandler)UrlDecode
```go
func (baseHandler *BaseHandler) UrlDecode(value string) string
```
```UrlDecode``` url decode.
### func (*BaseHandler)BeforeRequest<a name="BeforeRequest"></a>
```go
func (baseHandler *BaseHandler) BeforeRequest()
```
```BeforeRequest``` run this function before processing http request.
### func (*BaseHandler)TeardownRequest<a name="TeardownRequest"></a>
```go
func (baseHandler *BaseHandler) TeardownRequest()
```
```TeardownRequest```run this function when teardown http request.
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
### func (urlPattern *UrlPattern) AppendRouterPattern<a name="AppendRouterPattern"></a>
```go
func (urlPattern *UrlPattern) AppendRouterPattern(router Router, v interface {
	Handle(http.ResponseWriter, *http.Request)
})
```
Use this method to append a handler to an url and setup middleware to the handler.  
Example:
```go
// WithTracing get request uri
func WithTracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Tracing request for %s", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

var routers = []TigoWeb.Router{
	{Url: "/test", Handler: &TestHandler{}, Middleware: []TigoWeb.Middleware{WithTracing}},
}

func main() {
	application := TigoWeb.Application{IPAddress: "0.0.0.0", Port: 8080, UrlRouters: routers}
	application.Run()
}
```
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
### func (application *Application)Listen
```go
func (application *Application)Listen(port int)
```
Define the port for server to listen. **Attention**: The value of port in configuration will cover the value of port in this method.
### func (application *Application)MountFileServer
```go
func (application *Application)MountFileServer(dir string, uris ...string)
```
Mount http file server.
- `dir` the directory of folder will be mounted 
- `uris` the slice of uris
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
url encode<a name="urlencode"></a>
```go
func UrlEncode(value string)
```
Use this method to urlencode a string.  
url decode<a name="urldecode"></a>
```go
func UrlDecode(value string)
```
Use this method to urldecode a string.  
Initialize global configuration<a name="InitGlobalConfig"></a>
```go
func InitGlobalConfig(configPath string)
```
Use this method to initialize global configuration.  
Initialize global configuration<a name="InitGlobalConfigWithObj"></a>
```go
func InitGlobalConfigWithObj(config GlobalConfig)
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

func (helloHandler *HelloHandler)Get() {
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
    "error": "/Users/karllee/Desktop/run.log",
    "time_roll": "H*2"  // slice log file every 2 hours
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
//   - TimeRoll
// discard: discard, stdout: print in console; the path of log file
type LogLevel struct {
    Trace    string   `json:"trace"`
    Info     string   `json:"info"`
    Warning  string   `json:"warning"`
    Error    string   `json:"error"`
    TimeRoll string   `json:"time_roll"`
}
```
Initialize this structure and pass the instance to ```InitLoggerWithObject``` to init logger package.  
TimeRoll:
- D: slice log file by Day, E.g: "D*6" slice log file every six days.
- H: slice log file by Hour, E.g: "H*6" slice log file every six hours.
- M: slice log file by Minute, E.g: "M*6" slice log file every six minutes.
- S: slice log file by Second, E.g: "S*6" slice log file every six seconds.
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
# Tigo.request<a name="request"></a>
Tigo.request provides some functions to request an url.
## type Response<a name="httpResponse"></a>
```go
type Response struct {
    *http.Response
    Content []byte
}
```
HTTP Response Instance.
### func (response *Response)ToContentStr<a name="InitWarning"></a>
```go
func (response *Response)ToContentStr() string
```
Convert `Response.Content` to string type.
## request module<a name="requestFunctions"></a>
### func Request<a name="Request"></a>
```go
func Request(method string, requestUrl string, postParams map[string]interface{}, headers ...map[string]string) (*Response, error)
```
Make requestion to an url.
### func MakeRequest
```go
func MakeRequest(method string, requestUrl string, bodyReader io.Reader, headers ...map[string]string) (*Response, error)
```
Make requestion to an url.
### func Get<a name="Get"></a>
```go
func Get(requestUrl string, headers ...map[string]string) (*Response, error)
```
Http Get method.
### func Post<a name="Post"></a>
```go
func Post(requestUrl string, postParams map[string]interface{}, headers ...map[string]string) (*Response, error)
```
Http Post method.
### func Put<a name="Put"></a>
```go
func Put(requestUrl string, postParams map[string]interface{}, headers ...map[string]string) (*Response, error)
```
Http Put method.
### func Patch<a name="Patch"></a>
```go
func Patch(requestUrl string, postParams map[string]interface{}, headers ...map[string]string) (*Response, error)
```
Http Patch method.
### func Head<a name="Head"></a>
```go
func Head(requestUrl string, headers ...map[string]string) (*Response, error)
```
Http Head method.
### func Options<a name="Options"></a>
```go
func Options(requestUrl string, headers ...map[string]string) (*Response, error)
```
Http Options method.
### func Delete<a name="Delete"></a>
```go
func Delete(requestUrl string, headers ...map[string]string) (*Response, error)
```
Http Delete method.
# Tigo.binding<a name="binding"></a>
binding provides some functions to verify structure instance.
## binding functions<a name="bindingFunctions"></a>
### func ParseJsonToInstance<a name="ParseJsonToInstance"><a/>
```go
func ParseJsonToInstance(jsonBytes []byte, obj interface{}) error
```
Verify the instance unmarshal from the `jsonBytes`.
### func ValidateInstance<a name="ValidateInstance"><a/>
```go
func ValidateInstance(obj interface{}) error
```
Use this method to verify the `obj` with tag.
```go
type Company struct {
    Name string `json:"name" required:"false"`
    Addr string `json:"addr" required:"false"`
}

func (c *Company) Check() (err error) {
	if len(c.Name) > 100 {
		return errors.New("Company.Name is invalid")
	}
	return 
}

type Boss struct {
    Name    string  `json:"name" required:"true"`
    Age     int     `json:"age" required:"true" default:"18"`
    Company Company `json:"company" required:"true"`
}
/*This is OK👌*/

type Stuff struct {
    Name    string   `json:"name" required:"true"`
    Age     int      `json:"age" required:"true" default:"18"`
    Company *Company `json:"company" required:"true"`  // OK
}
/*This is OK👌*/

// Check method will run when checking `Stuff`.
func (s *Stuff) Check() (err error){
	if s.Age < 18 {
		return errors.New("Stuff.Age is invalid!")
	}
	return s.Company.Check()
}

type Others struct {
    Name    string   `json:"name" required:"true"`
    Age     *int     `json:"age" required:"true" default:"18"`  // Not Support
    Company Company  `json:"company" required:"true"`
}
/*This is OK👌*/
```
### func FormBytesToStructure<a name="FormBytesToStructure"><a/>
```go
func FormBytesToStructure(form []byte, obj interface{}) error
```
Use this method to convert form to structure instance.
### func ParseFormToInstance<a name="ParseFormToInstance"><a/>
```go
func ParseFormToInstance(form []byte, obj interface{}) error
```
Verify the instance unmarshal from the `form`.
