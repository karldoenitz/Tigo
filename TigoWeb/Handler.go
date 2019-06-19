// Package TigoWeb Copyright 2018 The Tigo Authors. All rights reserved.
package TigoWeb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/karldoenitz/Tigo/binding"
	"github.com/karldoenitz/Tigo/logger"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

// BaseHandler 是Handler的基础类，开发者开发的handler继承此类
type BaseHandler struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	JsonParams     map[string]interface{}
	ctxValMap      map[string]interface{}
}

// InitHandler 初始化Handler的方法
func (baseHandler *BaseHandler) InitHandler(responseWriter http.ResponseWriter, request *http.Request) {
	baseHandler.Request = request
	baseHandler.ResponseWriter = responseWriter
	baseHandler.Request.ParseForm()
	baseHandler.ctxValMap = map[string]interface{}{}
}

// GetBody 获取HTTP报文体
func (baseHandler *BaseHandler) GetBody() []byte {
	if _, ok := baseHandler.ctxValMap["body"]; !ok {
		body, err := ioutil.ReadAll(baseHandler.Request.Body)
		if err != nil {
			logger.Error.Println(err.Error())
			return nil
		}
		defer func() {
			ioReader := ioutil.NopCloser(bytes.NewBuffer(body))
			baseHandler.Request.Body = ioReader
		}()
		baseHandler.ctxValMap["body"] = body
	}
	return baseHandler.ctxValMap["body"].([]byte)
}

// PassJson 用来解析json中的值
func (baseHandler *BaseHandler) PassJson() {
	if baseHandler.GetHeader("Content-Type") == "application/json" {
		jsonData := baseHandler.GetBody()
		//使用 json.Unmarshal(data []byte, v interface{})进行转换，返回 error 信息
		err := json.Unmarshal(jsonData, &baseHandler.JsonParams)
		if err != nil {
			logger.Error.Println(err.Error())
		}
	}
}

/////////////////////////////////////////////////////output/////////////////////////////////////////////////////////////

// ToJson 将对象转化为Json字符串，转换失败则返回空字符串。
// 传入参数Response为一个interface，必须有成员函数Print。
func (baseHandler *BaseHandler) ToJson(response interface{}) (result string) {
	// 将该对象转换为byte字节数组
	jsonResult, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		logger.Error.Println(jsonErr.Error())
		return ""
	}
	// 将byte数组转换为string
	return string(jsonResult)
}

// ResponseAsJson 向客户端响应一个Json结果，默认字符集为utf-8
func (baseHandler *BaseHandler) ResponseAsJson(response interface{}, charset ...string) {
	// 将对象转换为Json字符串
	jsonResult := baseHandler.ToJson(response)
	// 设置http报文头内的Content-Type
	if len(charset) > 0 {
		baseHandler.ResponseWriter.Header().Set("Content-Type", fmt.Sprintf("application/json; %s", charset[0]))
	} else {
		baseHandler.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	fmt.Fprintf(baseHandler.ResponseWriter, jsonResult)
}

// ResponseAsText 向客户端响应一个Text结果
func (baseHandler *BaseHandler) ResponseAsText(result string) {
	fmt.Fprintf(baseHandler.ResponseWriter, result)
}

// ResponseAsHtml 向客户端响应一个html结果，默认字符集为utf-8
func (baseHandler *BaseHandler) ResponseAsHtml(result string, charset ...string) {
	if len(charset) > 0 {
		baseHandler.ResponseWriter.Header().Set("Content-Type", fmt.Sprintf("text/html; %s", charset[0]))
	} else {
		baseHandler.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	fmt.Fprintf(baseHandler.ResponseWriter, result)
}

// Response 向客户端响应一个结果
func (baseHandler *BaseHandler) Response(result ...interface{}) {
	fmt.Fprintf(baseHandler.ResponseWriter, "%v", result)
}

// ResponseFmt 向客户端响应一个字符串，支持format格式化字符串
func (baseHandler *BaseHandler) ResponseFmt(format string, values ...interface{}) {
	fmt.Fprintf(baseHandler.ResponseWriter, format, values...)
}

// ServerError 将服务器端发生的错误返回给客户端
func (baseHandler *BaseHandler) ServerError(err error) {
	http.Error(baseHandler.ResponseWriter, err.Error(), http.StatusInternalServerError)
}

// Render 渲染模板，返回数据
// 参数解析如下：
//   - data：表示传入的待渲染的数据
//   - templates：表示模板文件的路径，接受多个模板文件
func (baseHandler *BaseHandler) Render(data interface{}, templates ...string) {
	templateBasePath := globalConfig.Template
	var templatePath []string
	for _, value := range templates {
		value = templateBasePath + value
		templatePath = append(templatePath, value)
	}
	t, err := template.ParseFiles(templatePath...)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	er := t.Execute(baseHandler.ResponseWriter, data)
	if er != nil {
		logger.Error.Println(er.Error())
	}
}

// RedirectPermanently 向客户端永久重定向一个地址
func (baseHandler *BaseHandler) RedirectPermanently(url string) {
	baseHandler.ResponseWriter.WriteHeader(301)
	baseHandler.SetHeader("Location", url)
	fmt.Fprintf(baseHandler.ResponseWriter, "")
}

// Redirect 向客户端暂时重定向一个地址
func (baseHandler *BaseHandler) Redirect(url string, expire ...time.Time) {
	baseHandler.SetHeader("Location", url)
	baseHandler.ResponseWriter.WriteHeader(302)
	if len(expire) > 0 {
		expireTime := expire[0]
		expires := expireTime.Format("Mon, 02 Jan 2006 15:04:05 GMT")
		baseHandler.SetHeader("Expires", expires)
	}
	fmt.Fprintf(baseHandler.ResponseWriter, "")
}

/////////////////////////////////////////////////////cookie/////////////////////////////////////////////////////////////

// SetCookie 设置cookie
func (baseHandler *BaseHandler) SetCookie(name string, value string) {
	cookie := http.Cookie{Name: name, Value: value}
	http.SetCookie(baseHandler.ResponseWriter, &cookie)
}

// SetCookieObject 设置高级cookie选项
func (baseHandler *BaseHandler) SetCookieObject(cookie Cookie) {
	responseCookie := cookie.ToHttpCookie()
	http.SetCookie(baseHandler.ResponseWriter, &responseCookie)
}

// SetSecureCookie 设置加密cookie
func (baseHandler *BaseHandler) SetSecureCookie(name string, value string, key ...string) {
	securityKey := ""
	if len(key) > 0 {
		securityKey = key[0]
	} else {
		securityKey = globalConfig.Cookie
	}
	cookie := Cookie{
		Name:        name,
		Value:       value,
		IsSecurity:  true,
		SecurityKey: securityKey,
	}
	baseHandler.SetCookieObject(cookie)
}

// SetAdvancedCookie 设置cookie
//  - name cookie的名称
//  - value cookie的value
//  - attrs cookie的其他属性值，示例如下：
//    - "path={{string}}" 设置cookie的有效作用地址
//    - "domain={{string}}" 设置cookie的作用域
//    - "raw={{string}}" 设置cookie的raw值
//    - "maxAge={{int}}" 设置cookie的MaxAge，表示未指定“Max-Age”属性，表示现在删除cookie，相当于'Max-Age：0'，表示Max-Age属性存在并以秒为单位给出
//    - "expires={{int}}" 设置cookie的过期时间，按秒计算
//    - "secure={{bool}}" 设置cookie是否只限于加密传输
//    - "httpOnly={{bool}}" 设置cookie是否只限于http/https传输
//    - "isSecurity={{bool}}" 设置cookie是否要进行加密
func (baseHandler *BaseHandler) SetAdvancedCookie(name string, value string, attrs ...string) {
	key := globalConfig.Cookie
	var Path, Domain, Raw string
	var IsSecurity, Secure, HttpOnly bool
	var Expires time.Time
	var MaxAge int

	for _, attr := range attrs {
		switch {
		case strings.HasPrefix(attr, "path="):
			Path = strings.Replace(attr, "path=", "", 1)
			break
		case strings.HasPrefix(attr, "domain="):
			Domain = strings.Replace(attr, "domain=", "", 1)
			break
		case strings.HasPrefix(attr, "raw="):
			Raw = strings.Replace(attr, "raw=", "", 1)
			break
		case strings.HasPrefix(attr, "maxAge="):
			tmp := strings.Replace(attr, "maxAge=", "", 1)
			MaxAge, _ = strconv.Atoi(tmp)
			break
		case strings.HasPrefix(attr, "expires="):
			tmp := strings.Replace(attr, "expires=", "", 1)
			second, _ := strconv.Atoi(tmp)
			Expires = time.Now().Add(time.Second * time.Duration(second))
			break
		case attr == "secure=true":
			Secure = true
			break
		case attr == "httpOnly=true":
			HttpOnly = true
			break
		case attr == "isSecurity=true":
			IsSecurity = true
			break
		}
	}

	cookie := Cookie{
		Name:        name,
		Value:       value,
		IsSecurity:  IsSecurity,
		SecurityKey: key,
		Path:        Path,
		Domain:      Domain,
		Expires:     Expires,
		MaxAge:      MaxAge,
		Secure:      Secure,
		HttpOnly:    HttpOnly,
		Raw:         Raw,
	}
	baseHandler.SetCookieObject(cookie)
}

// GetCookie 获取cookie值，如果获取失败则返回空字符串
func (baseHandler *BaseHandler) GetCookie(name string) (value string) {
	cookie, err := baseHandler.Request.Cookie(name)
	if err != nil {
		return ""
	}
	value = cookie.Value
	return value
}

// GetSecureCookie 获取加密cookie值，如果获取失败则返回空
func (baseHandler *BaseHandler) GetSecureCookie(name string, key ...string) (value string) {
	securityKey := ""
	if len(key) > 0 {
		securityKey = key[0]
	} else {
		securityKey = globalConfig.Cookie
	}
	httpCookie, err := baseHandler.Request.Cookie(name)
	if err != nil {
		return ""
	}
	cookie := Cookie{}
	cookie.ConvertFromHttpCookie(*httpCookie)
	cookie.IsSecurity = true
	cookie.SecurityKey = securityKey
	value = cookie.GetCookieDecodeValue()
	return value
}

// GetCookieObject 获取cookie对象，多参数输入，参数如下：
//   - 无参数：默认cookieName为空字符串
//   - 一个参数：传入的参数为cookieName
//   - 多个参数：传入的第一个参数为cookieName，第二个参数为加密/解密cookie所用的Key，此时认为cookie是需要进行加密/解密处理的
func (baseHandler *BaseHandler) GetCookieObject(name ...string) (Cookie, error) {
	cookie := Cookie{}
	var cookieName, key string
	length := len(name)
	switch {
	case length < 1:
		cookieName = ""
	case length == 1:
		cookieName = name[0]
	case length > 1:
		cookieName = name[0]
		key = name[1]
	}
	httpCookie, err := baseHandler.Request.Cookie(cookieName)
	if err != nil {
		return cookie, nil
	}
	cookie.ConvertFromHttpCookie(*httpCookie)
	if len(key) > 0 {
		cookie.SetSecurityKey(key)
	}
	return cookie, nil
}

// ClearCookie 清除当前path下的指定的cookie
func (baseHandler *BaseHandler) ClearCookie(name string) {
	cookie := Cookie{
		Name:    name,
		Expires: time.Now(),
	}
	baseHandler.SetCookieObject(cookie)
}

// ClearAllCookie 清除当前path下所有的cookie
func (baseHandler *BaseHandler) ClearAllCookie() {
	cookies := baseHandler.Request.Cookies()
	for _, cookie := range cookies {
		baseHandler.ClearCookie(cookie.Name)
	}
}

/////////////////////////////////////////////////////input//////////////////////////////////////////////////////////////

// GetHeader 获取header
func (baseHandler *BaseHandler) GetHeader(name string) (value string) {
	value = baseHandler.Request.Header.Get(name)
	return value
}

// SetHeader 设置header
func (baseHandler *BaseHandler) SetHeader(name string, value string) {
	baseHandler.ResponseWriter.Header().Set(name, value)
}

// GetParameter 根据key获取对应的参数值
//   - 如果Content-Type是application/json，则直接从http的body中解析出key对应的value
//   - 否则，根据key直接获取value
func (baseHandler *BaseHandler) GetParameter(key string) (value *ReqParams) {
	jsonValue := &ReqParams{}
	if baseHandler.GetHeader("Content-Type") == "application/json" {
		if value, ok := baseHandler.JsonParams[key]; ok {
			jsonValue.Value = value
		} else {
			jsonValue.Value = nil
		}
		return jsonValue
	}
	jsonValue.Value = baseHandler.Request.FormValue(key)
	return jsonValue
}

// GetJsonValue 根据key获取对应的参数值，解析json数据，返回对应的value
func (baseHandler *BaseHandler) GetJsonValue(key string) interface{} {
	var mapResult map[string]interface{}
	jsonData := baseHandler.GetBody()
	//使用 json.Unmarshal(data []byte, v interface{})进行转换，返回 error 信息
	err := json.Unmarshal(jsonData, &mapResult)
	if err != nil {
		return ""
	}
	return mapResult[key]
}

// BeforeRequest 在每次响应HTTP请求之前执行此函数
func (baseHandler *BaseHandler) BeforeRequest() {
	return
}

// TeardownRequest 在每次响应HTTP请求之后执行此函数
func (baseHandler *BaseHandler) TeardownRequest() {
	return
}

//////////////////////////////////////////////////HTTP Method///////////////////////////////////////////////////////////

// 请求方法不合法
func (baseHandler *BaseHandler) methodNotAllowed() {
	baseHandler.ResponseWriter.WriteHeader(405)
}

// Get 方法
func (baseHandler *BaseHandler) Get() {
	baseHandler.methodNotAllowed()
}

// Put 方法
func (baseHandler *BaseHandler) Put() {
	baseHandler.methodNotAllowed()
}

// Post 方法
func (baseHandler *BaseHandler) Post() {
	baseHandler.methodNotAllowed()
}

// Connect 方法
func (baseHandler *BaseHandler) Connect() {
	baseHandler.methodNotAllowed()
}

// Trace 方法
func (baseHandler *BaseHandler) Trace() {
	baseHandler.methodNotAllowed()
}

// Head 方法
func (baseHandler *BaseHandler) Head() {
}

// Delete 方法
func (baseHandler *BaseHandler) Delete() {
	baseHandler.methodNotAllowed()
}

// Options 方法
func (baseHandler *BaseHandler) Options() {
	baseHandler.methodNotAllowed()
}

//////////////////////////////////////////////////Context Method////////////////////////////////////////////////////////

// SetCtxVal 在上下文中设置值
func (baseHandler *BaseHandler) SetCtxVal(key string, val interface{}) {
	baseHandler.ctxValMap[key] = val
}

// GetCtxVal 从上下文获取值
func (baseHandler *BaseHandler) GetCtxVal(key string) interface{} {
	if val, isExisted := baseHandler.ctxValMap[key]; isExisted {
		return val
	}
	return nil
}

//////////////////////////////////////////////////http message dump/////////////////////////////////////////////////////

// 获取http请求报文
func (baseHandler *BaseHandler) getHttpRequestMsg() string {
	req, err := httputil.DumpRequest(baseHandler.Request, true)
	if err != nil {
		return err.Error()
	}
	if baseHandler.Request.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		bodyData := getFormDataStr(baseHandler.Request.Form)
		return fmt.Sprintf("%s%s", string(req), bodyData)
	}
	return string(req)
}

// DumpHttpRequestMsg 获取http请求报文，根据logLevel值进行不同的输出
//  - 1: 将http报文输出到trace级别日志中
//  - 2: 将http报文输出到info级别日志中
//  - 3: 将http报文输出到warning级别日志中
//  - 4: 将http报文输出到error级别日志中
//  - others: 将http报文输出到控制台
func (baseHandler *BaseHandler) DumpHttpRequestMsg(logLevel int) {
	reqMsg := baseHandler.getHttpRequestMsg()
	switch logLevel {
	case logger.TraceLevel:
		logger.Trace.Println(reqMsg)
	case logger.InfoLevel:
		logger.Info.Println(reqMsg)
	case logger.WarningLevel:
		logger.Warning.Println(reqMsg)
	case logger.ErrorLevel:
		logger.Error.Println(reqMsg)
	default:
		fmt.Println(reqMsg)
	}
}

////////////////////////////////////////////////////////utils///////////////////////////////////////////////////////////

// CheckJsonBinding 检查提交的json是否符合要求
func (baseHandler *BaseHandler) CheckJsonBinding(obj interface{}) error {
	jsonData := baseHandler.GetBody()
	return binding.ParseJsonToInstance(jsonData, obj)
}

// CheckFormBinding 检查提交的form是否符合要求
func (baseHandler *BaseHandler) CheckFormBinding(obj interface{}) error {
	if err := binding.UnmarshalForm(baseHandler.Request.Form, obj); err != nil {
		return err
	}
	return binding.ValidateInstance(obj)
}

// CheckParamBinding 检查提交的参数是否符合要求
func (baseHandler *BaseHandler) CheckParamBinding(obj interface{}) error {
	if baseHandler.GetHeader("Content-Type") == "application/json" {
		return baseHandler.CheckJsonBinding(obj)
	}
	if baseHandler.GetHeader("Content-Type") == "application/x-www-form-urlencoded" {
		return baseHandler.CheckFormBinding(obj)
	}
	return nil
}

// CheckUrlParamBinding 检查url上传递的参数是否符合要求
func (baseHandler *BaseHandler) CheckUrlParamBinding(obj interface{}) error {
	return baseHandler.CheckFormBinding(obj)
}
