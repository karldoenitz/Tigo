// Package TigoWeb Copyright 2018 The Tigo Authors. All rights reserved.
package TigoWeb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/karldoenitz/Tigo/binding"
	"github.com/karldoenitz/Tigo/logger"
	"gorm.io/gorm"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"reflect"
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
	if err := baseHandler.Request.ParseForm(); err != nil {
		logger.Warning.Println(err.Error())
	}
	baseHandler.ctxValMap = map[string]interface{}{}
	baseHandler.JsonParams = map[string]interface{}{}
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
	if strings.Contains(baseHandler.GetHeader("Content-Type"), "application/json") {
		jsonData := baseHandler.GetBody()
		// 使用 json.Unmarshal(data []byte, v interface{})进行转换，返回 error 信息
		err := json.Unmarshal(jsonData, &baseHandler.JsonParams)
		if err != nil {
			logger.Error.Println(err.Error())
		}
	}
}

/////////////////////////////////////////////////////output/////////////////////////////////////////////////////////////

// ToJson 将对象转化为Json字符串，转换失败则返回空字符串。
// 传入参数Response为一个interface，必须有成员函数Print。
//   - response: 需要转换成json的实例
func (baseHandler *BaseHandler) ToJson(response interface{}) (result []byte) {
	// 将该对象转换为byte字节数组
	jsonResult, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		logger.Error.Println(jsonErr.Error())
	}
	return jsonResult
}

// ToJsonStr 将对象转化为Json字符串，转换失败则返回空字符串。
// 传入参数Response为一个interface，必须有成员函数Print。
//   - response: 需要转换成json的实例
func (baseHandler *BaseHandler) ToJsonStr(response interface{}) (result string) {
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
//   - response: 需要响应给客户端的数据
//   - charset: 数据集，默认utf-8编码
func (baseHandler *BaseHandler) ResponseAsJson(response interface{}, charset ...string) {
	// 将对象转换为Json字符串
	jsonResult := baseHandler.ToJson(response)
	// 设置http报文头内的Content-Type
	if len(charset) > 0 {
		baseHandler.ResponseWriter.Header().Set("Content-Type", fmt.Sprintf("application/json; %s", charset[0]))
	} else {
		baseHandler.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	if _, err := baseHandler.ResponseWriter.Write(jsonResult); err != nil {
		logger.Warning.Println(err.Error())
	}
}

// ResponseAsText 向客户端响应一个Text结果
//   - response: 需要返回的文本内容
func (baseHandler *BaseHandler) ResponseAsText(result string) {
	if _, err := fmt.Fprintf(baseHandler.ResponseWriter, "%s", result); err != nil {
		logger.Warning.Println(err.Error())
	}
}

// ResponseAsHtml 向客户端响应一个html结果，默认字符集为utf-8
//   - result: 相应的结果
//   - charset: 字符集，默认utf-8
func (baseHandler *BaseHandler) ResponseAsHtml(result string, charset ...string) {
	cs := "text/html; charset=utf-8"
	if len(charset) > 0 {
		cs = fmt.Sprintf("text/html; %s", charset[0])
	}
	baseHandler.ResponseWriter.Header().Set("Content-Type", cs)
	if _, err := fmt.Fprintf(baseHandler.ResponseWriter, "%s", result); err != nil {
		logger.Warning.Println(err.Error())
	}
}

// Response 向客户端响应一个结果
func (baseHandler *BaseHandler) Response(result ...interface{}) {
	if _, err := fmt.Fprintf(baseHandler.ResponseWriter, "%v", result); err != nil {
		logger.Warning.Println(err.Error())
	}
}

// ResponseFmt 向客户端响应一个字符串，支持format格式化字符串
//   - format: 格式化母串
//   - values: 返回的值
func (baseHandler *BaseHandler) ResponseFmt(format string, values ...interface{}) {
	if _, err := fmt.Fprintf(baseHandler.ResponseWriter, format, values...); err != nil {
		logger.Warning.Println(err.Error())
	}
}

// ResponseWithFilter 通过filter返回请求结果，目前只会解析url上传的参数，当前只支持Get请求。
//   - filter: 过滤器对象，不要传指针和引用
//   - conn: 数据库的链接
func (baseHandler *BaseHandler) ResponseWithFilter(filter interface{}, conn *gorm.DB, model interface{}) {
	conn = conn.Model(model)
	filterType := reflect.TypeOf(filter)
	var offset, limit, params string
	url := strings.Split(baseHandler.Request.URL.String(), "?")
	if len(url) > 1 {
		params = "&" + url[1]
	}
	for i := 0; i < filterType.NumField(); i++ {
		field := filterType.Field(i)
		urlParam := GetTagValue(field, "url")
		columnName := GetTagValue(field, "column")
		// 判断 url 上是否有这个形参，没有形参则忽略
		if !strings.Contains(params, fmt.Sprintf("&%s=", urlParam)) {
			continue
		}
		// 有形参则取值
		value := baseHandler.Request.FormValue(urlParam)
		// TODO 这里分页后续优化一下
		if columnName == "offset" {
			offset = value
			continue
		}
		if columnName == "limit" {
			limit = value
			continue
		}
		conn = convertCondition(urlParam, columnName, value, conn)
	}
	var total, cnt int64
	conn.Count(&total)
	if offset != "" {
		conn = convertCondition("offset", "offset", offset, conn)
	}
	if limit != "" {
		conn = convertCondition("limit", "limit", limit, conn)
	}
	conn.Count(&cnt)
	t := reflect.TypeOf(model)
	d := reflect.New(reflect.MakeSlice(reflect.SliceOf(t), 0, 0).Type()).Interface()
	conn.Find(d)
	result := reflect.ValueOf(filter).MethodByName("Process").Call([]reflect.Value{reflect.ValueOf(d)})
	baseHandler.ResponseAsJson(map[string]interface{}{
		"total": total,
		"count": cnt,
		"data":  result[0].Interface(),
	})
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

// redirectPermanently 向客户端永久重定向一个地址
func (baseHandler *BaseHandler) redirectPermanently(url string, status int) {
	baseHandler.SetHeader("Location", url)
	baseHandler.ResponseWriter.WriteHeader(status)
	if _, err := baseHandler.ResponseWriter.Write(nil); err != nil {
		logger.Warning.Println(err.Error())
	}
}

// redirect 向客户端暂时重定向一个地址
func (baseHandler *BaseHandler) redirect(url string, status int, expire ...time.Time) {
	baseHandler.SetHeader("Location", url)
	baseHandler.ResponseWriter.WriteHeader(status)
	if len(expire) > 0 {
		expireTime := expire[0]
		expires := expireTime.Format("Mon, 02 Jan 2006 15:04:05 GMT")
		baseHandler.SetHeader("Expires", expires)
	}
	if _, err := baseHandler.ResponseWriter.Write(nil); err != nil {
		logger.Warning.Println(err.Error())
	}
}

// MovePermanently 向客户端永久性移动一个地址
//   - url: 指定客户端要移动的目标地址
func (baseHandler *BaseHandler) MovePermanently(url string) {
	baseHandler.redirectPermanently(url, http.StatusMovedPermanently)
}

// Move 向客户端暂时移动一个地址
//   - url: 指定客户端要移动的目标地址
//   - expire: 过期时间
func (baseHandler *BaseHandler) Move(url string, expire ...time.Time) {
	baseHandler.redirect(url, http.StatusFound, expire...)
}

// RedirectPermanently 向客户端永久重定向一个地址
//   - url: 重定向的地址
func (baseHandler *BaseHandler) RedirectPermanently(url string) {
	baseHandler.redirectPermanently(url, http.StatusPermanentRedirect)
}

// Redirect 向客户端暂时重定向一个地址
//   - url: 重定向的地址
//   - expire: 过期时间
func (baseHandler *BaseHandler) Redirect(url string, expire ...time.Time) {
	baseHandler.redirect(url, http.StatusTemporaryRedirect, expire...)
}

// RedirectTo 自定义重定向
//   - url: 重定向的url
//   - status: http状态码
//   - expire: 过期时间
func (baseHandler *BaseHandler) RedirectTo(url string, status int, expire ...time.Time) {
	baseHandler.redirect(url, status, expire...)
}

/////////////////////////////////////////////////////cookie/////////////////////////////////////////////////////////////

// SetCookie 设置cookie
// SetCookie未设置cookie的domain及path，此cookie仅对当前路径有效，设置其他路径cookie可参考SetAdvancedCookie
//   - name: cookie的name
//   - value: cookie的值
func (baseHandler *BaseHandler) SetCookie(name string, value string) {
	cookie := http.Cookie{Name: name, Value: value}
	http.SetCookie(baseHandler.ResponseWriter, &cookie)
}

// SetCookieObject 设置高级cookie选项
//   - cookie: 要设置的cookie对象
func (baseHandler *BaseHandler) SetCookieObject(cookie Cookie) {
	responseCookie := cookie.ToHttpCookie()
	http.SetCookie(baseHandler.ResponseWriter, &responseCookie)
}

// SetSecureCookie 设置加密cookie
// SetSecureCookie未设置cookie的domain及path，此cookie仅对当前路径有效，设置其他路径cookie可参考SetAdvancedCookie
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
//   - name cookie的名称
//   - value cookie的value
//   - attrs cookie的其他属性值，示例如下：
//   - "path={{string}}" 设置cookie的有效作用地址
//   - "domain={{string}}" 设置cookie的作用域
//   - "raw={{string}}" 设置cookie的raw值
//   - "maxAge={{int}}" 设置cookie的MaxAge，表示未指定“Max-Age”属性，表示现在删除cookie，相当于'Max-Age：0'，表示Max-Age属性存在并以秒为单位给出
//   - "expires={{int}}" 设置cookie的过期时间，按秒计算
//   - "secure={{bool}}" 设置cookie是否只限于加密传输
//   - "httpOnly={{bool}}" 设置cookie是否只限于http/https传输
//   - "isSecurity={{bool}}" 设置cookie是否要进行加密
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
			Expires = time.Now().Local().Add(time.Second * time.Duration(second))
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
//   - name: cookie的name
func (baseHandler *BaseHandler) GetCookie(name string) (value string) {
	cookie, err := baseHandler.Request.Cookie(name)
	if err != nil {
		return ""
	}
	value = cookie.Value
	return value
}

// GetSecureCookie 获取加密cookie值，如果获取失败则返回空
//   - name: cookie的name
//   - key: cookie加密用的key
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

// ClearCookie 清除本次请求当前path下的指定的cookie
//   - name: 需要清楚的cookie的name
func (baseHandler *BaseHandler) ClearCookie(name string) {
	cookie := Cookie{
		Name:    name,
		Expires: time.Unix(0, 0),
		MaxAge:  0,
	}
	baseHandler.SetCookieObject(cookie)
}

// ClearAllCookie 清除本次请求当前path下所有的cookie
func (baseHandler *BaseHandler) ClearAllCookie() {
	cookies := baseHandler.Request.Cookies()
	for _, cookie := range cookies {
		cookie.Expires = time.Unix(0, 0)
		http.SetCookie(baseHandler.ResponseWriter, cookie)
	}
}

/////////////////////////////////////////////////////session////////////////////////////////////////////////////////////

// SetSession 根据key设置session值
//   - key: session对应的键
//   - value: session的值
func (baseHandler *BaseHandler) SetSession(key string, value interface{}) (err error) {
	sessionId := baseHandler.GetCookie(SessionCookieName)
	var session Session
	if sessionId == "" {
		// 此处先默认3600秒，下个版本改为从配置文件读取
		session = GlobalSessionManager.GenerateSession(0)
		sessionId = session.SessionId()
	} else {
		session = GlobalSessionManager.GetSessionBySid(sessionId)
		if session.SessionId() == "" {
			session = GlobalSessionManager.GenerateSession(0)
			sessionId = session.SessionId()
		}
	}
	baseHandler.SetAdvancedCookie(SessionCookieName, sessionId, "path=/")
	err = session.Set(key, value)
	return
}

// GetSession 根据key获取session值
func (baseHandler *BaseHandler) GetSession(key string, value interface{}) (err error) {
	sessionId := baseHandler.GetCookie(SessionCookieName)
	if sessionId == "" {
		logger.Info.Println("session id is empty")
		value = nil
		return
	}
	session := GlobalSessionManager.GetSessionBySid(sessionId)
	if session == nil {
		logger.Info.Println("session is nil")
		baseHandler.SetAdvancedCookie(SessionCookieName, "", "maxAge=0", "path=/")
		value = nil
		return
	}
	err = session.Get(key, value)
	return
}

// ClearSession 根据key清除对应的session值
func (baseHandler *BaseHandler) ClearSession(key string) {
	sessionId := baseHandler.GetCookie(SessionCookieName)
	if sessionId == "" {
		logger.Info.Println("session id is empty")
		return
	}
	session := GlobalSessionManager.GetSessionBySid(sessionId)
	if session == nil {
		logger.Info.Println("session is nil")
		baseHandler.SetAdvancedCookie(SessionCookieName, "", "maxAge=0", "path=/")
		return
	}
	session.Delete(key)
}

// DelSession 删除所有的session值
func (baseHandler *BaseHandler) DelSession() {
	sessionId := baseHandler.GetCookie(SessionCookieName)
	GlobalSessionManager.DeleteSession(sessionId)
	baseHandler.SetAdvancedCookie(SessionCookieName, "", "maxAge=0", "path=/")
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
	val := baseHandler.Request.FormValue(key)
	contentType := baseHandler.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if value, ok := baseHandler.JsonParams[key]; ok {
			jsonValue.Value = value
		} else if val != "" {
			jsonValue.Value = val
		} else {
			jsonValue.Value = nil
		}
		return jsonValue
	}
	jsonValue.Value = val
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

// GetPathParam 根据key获取Url上的参数
func (baseHandler *BaseHandler) GetPathParam(key string) (value PathParam) {
	vars := mux.Vars(baseHandler.Request)
	value = PathParam(vars[key])
	return
}

// GetPathParamStr 根据key获取Url上的参数
func (baseHandler *BaseHandler) GetPathParamStr(key string) string {
	return mux.Vars(baseHandler.Request)[key]
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
//   - demo请查看源代码中的注释
func (baseHandler *BaseHandler) SetCtxVal(key string, val interface{}) {
	/*
		在中间件中，如果初始化了一个handler，需要将handler内的Request和ResponseWriter作为参数传入到next.ServeHTTP中

		func Authorize(w *http.ResponseWriter, r *http.Request) bool {
			handler := TigoWeb.BaseHandler{Request: r, ResponseWriter: w}
			handler.SetCtxVal("key", "value")
			return true
		}
	*/
	ctx := baseHandler.Request.Context()
	ctx = context.WithValue(ctx, key, val)
	baseHandler.Request = baseHandler.Request.WithContext(ctx)
}

// GetCtxVal 从上下文获取值
func (baseHandler *BaseHandler) GetCtxVal(key string) interface{} {
	return baseHandler.Request.Context().Value(key)
}

//////////////////////////////////////////////////http message dump/////////////////////////////////////////////////////

// 获取http请求报文 TODO 校验此处是否正常
func (baseHandler *BaseHandler) getHttpRequestMsg() string {
	req, err := httputil.DumpRequest(baseHandler.Request, true)
	if err != nil {
		return err.Error()
	}
	if strings.Contains(baseHandler.GetHeader("Content-Type"), "application/x-www-form-urlencoded") {
		bodyData := getFormDataStr(baseHandler.Request.Form)
		return fmt.Sprintf("%s%s", string(req), bodyData)
	}
	return string(req)
}

// DumpHttpRequestMsg 获取http请求报文，根据logLevel值进行不同的输出
//   - 1: 将http报文输出到trace级别日志中
//   - 2: 将http报文输出到info级别日志中
//   - 3: 将http报文输出到warning级别日志中
//   - 4: 将http报文输出到error级别日志中
//   - others: 将http报文输出到控制台
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
	contentType := baseHandler.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
		return baseHandler.CheckJsonBinding(obj)
	}
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return baseHandler.CheckFormBinding(obj)
	}
	return nil
}

// CheckUrlParamBinding 检查url上传递的参数是否符合要求
func (baseHandler *BaseHandler) CheckUrlParamBinding(obj interface{}) error {
	return baseHandler.CheckFormBinding(obj)
}

// UrlEncode 对值进行url编码
func (baseHandler *BaseHandler) UrlEncode(value string) string {
	return UrlEncode(value)
}

// UrlDecode 对值进行url解码
func (baseHandler *BaseHandler) UrlDecode(value string) string {
	return UrlDecode(value)
}
