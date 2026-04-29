# Tigo Web 包架构文档

本文档详细分析 Tigo 框架 `web` 包的代码架构、设计模式以及每个模块的实现逻辑。

## 目录

1. [概述](#概述)
2. [架构设计](#架构设计)
3. [模块详解](#模块详解)
   - [Application (application.go)](#application-applicationgo)
   - [Handler (handler.go)](#handler-handlergo)
   - [URL路由 (urlpattern.go)](#url路由-urlpatterngo)
   - [中间件 (middleware.go)](#中间件-middlewarego)
   - [WebSocket支持](#websocket支持)
   - [配置管理 (struct.go, global.go)](#配置管理-structgo-globalgo)
   - [Session管理 (session.go)](#session管理-sessiongo)
   - [工具函数 (utils.go)](#工具函数-utilsgo)
4. [请求处理生命周期](#请求处理生命周期)
5. [README.md Demo 分析](#readmemd-demo-分析)

---

## 概述

`web` 包是 Tigo 框架的核心包，提供了搭建 Web 服务的所有基础功能。其核心设计理念是：

- **Handler 继承模式**：开发者通过继承 `BaseHandler` 实现业务逻辑
- **反射驱动的方法分发**：根据 HTTP 方法自动调用对应的 Handler 方法
- **中间件链模式**：支持灵活的请求预处理和后处理
- **插件化架构**：Session、WebSocket 等功能支持第三方扩展

### 包结构

```
web/
├── application.go        # 应用入口，服务启动
├── handler.go            # BaseHandler 实现
├── urlpattern.go         # URL 路由配置
├── urlpatternhandler.go  # HTTP Handler 路由处理
├── urlpatternwshandler.go # WebSocket Handler 路由处理
├── middleware.go         # 中间件实现
├── session.go            # Session 接口定义
├── struct.go             # 数据结构定义
├── global.go             # 全局配置变量
├── enums.go              # 常量和枚举定义
├── interface.go          # 核心接口定义
└── utils.go              # 工具函数
```

---

## 架构设计

### 核心架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        Application                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   Config    │  │  UrlPatterns │  │      mux.Router         │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Middleware Chain                             │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐    │
│  │ LogMiddleware  │→│ ErrorMiddleware │→│ CustomMiddleware │→...│
│  └────────────────┘  └────────────────┘  └────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Handler Lifecycle                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ InitHandler │→│BeforeRequest│→│ HTTP Method  │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
│         │                                      │                 │
│         ▼                                      ▼                 │
│  ┌─────────────┐                      ┌─────────────────┐       │
│  │  PassJson   │                      │TeardownRequest │       │
│  └─────────────┘                      └─────────────────┘       │
└─────────────────────────────────────────────────────────────────┘
```

### 设计模式

1. **模板方法模式 (Template Method)**
   - `BaseHandler` 定义了请求处理的骨架流程
   - 子类通过重写 `Get()`, `Post()` 等方法实现具体业务

2. **责任链模式 (Chain of Responsibility)**
   - 中间件通过 `chainMiddleware` 函数链接
   - 每个中间件可以决定是否继续执行下一个

3. **反射工厂模式**
   - `NewOriginHandler` 根据类型自动创建对应的 Handler
   - 通过反射检测 Handler 类型（普通/WebSocket）

4. **依赖注入模式**
   - `SessionInterface` 允许注入第三方 Session 实现
   - `WSConfig` 支持 WebSocket 配置注入

---

## 模块详解

### Application (application.go)

`Application` 是 Tigo 框架的入口点，负责初始化和启动 HTTP 服务。

#### 核心结构

```go
type Application struct {
    IPAddress   string      // IP 地址
    Port        int         // 端口
    UrlPatterns []Pattern   // URL 路由配置
    ConfigPath  string      // 配置文件路径
    muxRouter   *mux.Router // Gorilla Mux 路由器
}
```

#### 关键方法

| 方法 | 功能 |
|------|------|
| `Run()` | 启动 HTTP/HTTPS 服务 |
| `InitApp()` | 初始化配置和路由 |
| `Listen(port)` | 设置监听端口 |
| `StartSession()` | 启用 Session 功能 |
| `MountFileServer()` | 挂载静态文件服务 |
| `EndlessStart()` | 使用 endless 实现平滑重启 |
| `OverseerStart()` | 使用 overseer 实现平滑重启 |

#### 启动流程

```
Run()
  │
  ├── muxRouter = mux.NewRouter()
  │
  ├── InitApp()
  │     │
  │     ├── InitGlobalConfig() (如果配置了 ConfigPath)
  │     │
  │     └── UrlPattern.Init() → 挂载所有路由
  │
  └── run() → 启动 HTTP/HTTPS 服务
```

#### HTTPS 支持

当配置文件中设置了 `Cert` 和 `CertKey` 时，自动启用 HTTPS：

```go
case globalConfig != nil && globalConfig.Cert != "" && globalConfig.CertKey != "":
    httpServerErr = http.ListenAndServeTLS(address, globalConfig.Cert, globalConfig.CertKey, application.muxRouter)
```

---

### Handler (handler.go)

`BaseHandler` 是所有业务 Handler 的基类，提供了处理 HTTP 请求的核心功能。

#### 核心结构

```go
type BaseHandler struct {
    ResponseWriter http.ResponseWriter
    Request        *http.Request
    JsonParams     map[string]interface{}  // JSON 参数缓存
    ctxValMap      map[string]interface{}  // 上下文值缓存
}
```

#### 功能分类

##### 1. 初始化方法

```go
func (h *BaseHandler) InitHandler(w http.ResponseWriter, r *http.Request)
func (h *BaseHandler) PassJson()  // 解析 JSON 请求体
```

##### 2. 响应方法

| 方法 | 功能 |
|------|------|
| `ResponseAsJson(data)` | 返回 JSON 响应 |
| `ResponseAsText(text)` | 返回文本响应 |
| `ResponseAsHtml(html)` | 返回 HTML 响应 |
| `Render(data, templates...)` | 渲染模板文件 |
| `ResponseWithFilter(...)` | 带过滤器的 GORM 查询响应 |

##### 3. Cookie 方法

| 方法 | 功能 |
|------|------|
| `SetCookie(name, value)` | 设置普通 Cookie |
| `SetSecureCookie(name, value)` | 设置加密 Cookie |
| `SetAdvancedCookie(name, value, attrs...)` | 设置高级 Cookie |
| `GetCookie(name)` | 获取 Cookie 值 |
| `GetSecureCookie(name)` | 获取加密 Cookie 值 |
| `ClearCookie(name)` | 清除指定 Cookie |

##### 4. Session 方法

| 方法 | 功能 |
|------|------|
| `SetSession(key, value)` | 设置 Session |
| `GetSession(key, &value)` | 获取 Session |
| `ClearSession(key)` | 清除指定 Session |
| `DelSession()` | 删除所有 Session |

##### 5. 参数获取方法

| 方法 | 功能 |
|------|------|
| `GetParameter(key)` | 获取请求参数（支持 JSON/Form） |
| `GetJsonValue(key)` | 从 JSON body 获取值 |
| `GetPathParam(key)` | 获取 URL 路径参数 |
| `GetHeader(name)` | 获取请求头 |
| `GetBody()` | 获取请求体 |

##### 6. 重定向方法

| 方法 | 功能 |
|------|------|
| `Redirect(url)` | 临时重定向 |
| `RedirectPermanently(url)` | 永久重定向 |
| `Move(url)` | 临时移动 |
| `MovePermanently(url)` | 永久移动 |

##### 7. 生命周期钩子

```go
func (h *BaseHandler) BeforeRequest()   // 请求前钩子
func (h *BaseHandler) TeardownRequest() // 请求后钩子
```

##### 8. 数据绑定验证

```go
func (h *BaseHandler) CheckJsonBinding(obj interface{}) error
func (h *BaseHandler) CheckFormBinding(obj interface{}) error
func (h *BaseHandler) CheckParamBinding(obj interface{}) error
```

#### HTTP 方法处理

`BaseHandler` 为每个 HTTP 方法定义了默认实现，未实现的方法返回 405：

```go
func (h *BaseHandler) Get()     { h.methodNotAllowed() }
func (h *BaseHandler) Post()    { h.methodNotAllowed() }
func (h *BaseHandler) Put()     { h.methodNotAllowed() }
func (h *BaseHandler) Delete()  { h.methodNotAllowed() }
func (h *BaseHandler) Options() { h.methodNotAllowed() }
func (h *BaseHandler) Head()    { } // Head 方法默认不做任何处理
```

---

### URL路由 (urlpattern.go)

URL 路由模块负责将 URL 映射到对应的 Handler。

#### 核心结构

```go
// Pattern 路由对象
type Pattern struct {
    Url        string
    Handler    interface{}      // 可以是 Handler 类型或文件路径字符串
    Middleware []Middleware
}

// UrlPattern URL 路由集合
type UrlPattern struct {
    UrlPatterns []Pattern
    router      *mux.Router
}
```

#### 路由挂载流程

```go
func (urlPattern *UrlPattern) AppendRouterPattern(pattern Pattern, v OriginHandlerInterface) {
    // 1. 判断是否是文件服务器
    if filePath, isFileServer := pattern.Handler.(string); isFileServer {
        // 挂载文件服务器
        fileServer := http.FileServer(http.Dir(filePath))
        // ...
        return
    }

    // 2. 判断是否是 WebSocket Handler
    if _, isWSHandler := v.(*WSPatternHandle); isWSHandler {
        // 挂载 WebSocket 路由
        // ...
        return
    }

    // 3. 普通 HTTP Handler
    baseMiddleware := []middleware{HttpContextLogMiddleware, InternalServerErrorMiddleware}
    // 添加自定义中间件
    for _, mv := range pattern.Middleware {
        m := convertHandleFuncMV(mv)
        baseMiddleware = append(baseMiddleware, m)
    }
    middlewares := chainMiddleware(baseMiddleware...)
    urlPattern.router.HandleFunc(pattern.Url, middlewares(v.Handle))
}
```

#### Handler 类型检测

通过反射检测 Handler 类型：

```go
func NewOriginHandler(pattern Pattern) OriginHandlerInterface {
    handlerType := reflect.TypeOf(pattern.Handler)
    if handlerType.Kind() == reflect.Ptr {
        handlerType = handlerType.Elem()
    }
    // 检查是否有 WSBaseHandler 字段
    _, isWSHandler := handlerType.FieldByName(FnWebSocketBaseHandler)
    if isWSHandler {
        return newWSPatternHandle(handlerType, pattern)
    }
    return newUrlPatternHandle(handlerType, pattern)
}
```

---

### 中间件 (middleware.go)

中间件模块实现了请求预处理和后处理功能。

#### 中间件类型定义

```go
// 内部中间件类型
type middleware func(next http.HandlerFunc) http.HandlerFunc

// 用户定义的中间件类型
type Middleware func(*http.ResponseWriter, *http.Request) bool
```

#### 中间件链

```go
func chainMiddleware(mw ...middleware) middleware {
    return func(final http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            last := final
            for i := len(mw) - 1; i >= 0; i-- {
                last = mw[i](last)
            }
            last(w, r)
        }
    }
}
```

中间件链执行顺序：
```
请求 → Middleware1 → Middleware2 → ... → Handler → 响应
```

#### 内置中间件

##### 1. InternalServerErrorMiddleware

捕获 panic 并返回 500 错误：

```go
func InternalServerErrorMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if r := recover(); r != nil {
                // 处理 panic
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    }
}
```

##### 2. HttpContextLogMiddleware

记录请求日志和耗时：

```go
func HttpContextLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        startTime := time.Now()
        httpResponseWriter := HttpResponseWriter{w, http.StatusOK}
        defer func() {
            status := httpResponseWriter.GetStatus()
            duration := time.Now().Sub(startTime).Seconds() * 1e3
            logger.Info.Printf("%d | %fms | %s %s", status, duration, r.Method, r.RequestURI)
        }()
        next.ServeHTTP(&httpResponseWriter, r)
    }
}
```

---

### WebSocket支持

WebSocket 模块提供了完整的 WebSocket 功能支持。

#### WSConfig 配置

```go
type WSConfig struct {
    HandshakeTimeout  time.Duration            // 握手超时
    ReadBufferSize    int                      // 读缓冲区大小
    WriteBufferSize   int                      // 写缓冲区大小
    EnableCompression bool                     // 是否启用压缩
    PingInterval      time.Duration            // Ping 间隔
    PongWait          time.Duration            // Pong 等待超时
    WriteWait         time.Duration            // 写入超时
    CheckOrigin       func(*http.Request) bool // Origin 检查
}
```

#### WSBaseHandler 基类

```go
type WSBaseHandler struct {
    conn      *websocket.Conn
    ctx       context.Context
    cancelCtx context.CancelFunc
    writeMu   sync.Mutex      // 写锁，保证线程安全
    request   *http.Request
    config    WSConfig
}
```

#### 生命周期钩子

```go
func (h *WSBaseHandler) OnConnect()                    // 连接建立
func (h *WSBaseHandler) OnMessage(messageType, data)   // 收到消息
func (h *WSBaseHandler) OnError(err)                   // 发生错误
func (h *WSBaseHandler) OnDisconnect()                 // 连接关闭
```

#### 通信方法

| 方法 | 功能 |
|------|------|
| `SendText(msg)` | 发送文本消息（线程安全） |
| `SendJSON(v)` | 发送 JSON 消息（线程安全） |
| `SendBinary(data)` | 发送二进制消息（线程安全） |
| `ReadMessage()` | 读取消息 |
| `ReadJSON(v)` | 读取 JSON 消息 |
| `StartPingPong()` | 启动心跳 |

#### 默认通信流程

```go
func (h *WSBaseHandler) Communicate() {
    h.OnConnect()              // 1. 连接建立钩子
    for {
        messageType, data, err := h.conn.ReadMessage()
        if err != nil {
            h.OnError(err)
            break
        }
        h.OnMessage(messageType, data)  // 2. 消息处理
    }
    h.OnDisconnect()           // 3. 连接关闭钩子
}
```

#### WSPatternHandle 路由处理

```go
func (wsHandler *WSPatternHandle) Handle(w http.ResponseWriter, r *http.Request) {
    // 1. 创建 upgrader
    upgrader := &websocket.Upgrader{...}

    // 2. 升级 HTTP 连接
    conn, err := upgrader.Upgrade(w, r, nil)

    // 3. 创建 Handler 实例并注入依赖
    handler := reflect.New(wsHandler.Handler)
    // 设置 conn, request, config

    // 4. 调用 Communicate 方法
    VoidFuncCall(handler, FnWebSocketCommunicate)
}
```

---

### 配置管理 (struct.go, global.go)

#### GlobalConfig 全局配置

```go
type GlobalConfig struct {
    IP        string          `json:"ip"`        // IP 地址
    Port      int             `json:"port"`      // 端口
    Cert      string          `json:"cert"`      // HTTPS 证书路径
    CertKey   string          `json:"cert_key"`  // HTTPS 密钥路径
    Cookie    string          `json:"cookie"`    // Cookie 加密密钥
    Template  string          `json:"template"`  // 模板文件路径
    Log       logger.LogLevel `json:"log"`       // 日志配置
    WebSocket WSConfig        `json:"websocket"` // WebSocket 配置
}
```

#### 配置初始化

支持 JSON 和 YAML 两种格式：

```go
func (gc *GlobalConfig) Init(configPath string) {
    if strings.HasSuffix(configPath, ".json") {
        gc.initWithJson(configPath)
    }
    if strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, "yml") {
        gc.initWithYaml(configPath)
    }
}
```

#### Cookie 结构

```go
type Cookie struct {
    Name        string
    Value       string
    IsSecurity  bool      // 是否加密
    SecurityKey string    // 加密密钥
    Path        string
    Domain      string
    Expires     time.Time
    MaxAge      int
    Secure      bool
    HttpOnly    bool
    // ...
}
```

Cookie 加密/解密：

```go
func (c *Cookie) GetCookieEncodeValue() string  // 获取加密值
func (c *Cookie) GetCookieDecodeValue() string  // 获取解密值
```

---

### Session管理 (session.go)

Session 模块采用插件化设计，允许第三方实现。

#### 接口定义

```go
// SessionInterface - 插件需要实现此接口
type SessionInterface interface {
    NewSessionManager() SessionManager
}

// SessionManager - Session 管理器接口
type SessionManager interface {
    GenerateSession(expire int) Session
    GetSessionBySid(sid string) Session
    DeleteSession(sid string)
}

// Session - Session 操作接口
type Session interface {
    Set(key string, value interface{}) error
    Get(key string, value interface{}) error
    Delete(key string)
    SessionId() string
}
```

#### 全局变量

```go
var GlobalSessionManager SessionManager  // 全局 Session 管理器
var SessionCookieName = "TigoSessionId"  // Session Cookie 名称
```

#### 使用示例

```go
// 启用 Session（使用 tission 插件）
application.StartSession(tission.NewTission(), "MySessionId")

// 在 Handler 中使用
func (h *MyHandler) Get() {
    h.SetSession("user", "alice")
    user, _ := h.GetSession("user", &value)
}
```

---

### 工具函数 (utils.go)

#### 加密解密

```go
func Encrypt(src []byte, key []byte) (string, error)  // AES 加密
func Decrypt(src []byte, key []byte) ([]byte, error)  // AES 解密
func MD5(origin string) string                        // MD5 哈希
func MD5m16(origin string) string                     // 16位 MD5
```

加密流程：
```
原始数据 → AES-GCM 加密 → Base64 编码 → 加密结果
```

#### 反射工具

```go
// 调用无返回值方法
func VoidFuncCall(instance reflect.Value, funcName string, params ...reflect.Value)

// 获取结构体标签值
func GetTagValue(field reflect.StructField, tagKey string) string
```

#### GORM 查询构建

```go
func convertCondition(urlParam, column, value string, db *gorm.DB) *gorm.DB
```

支持的条件后缀：
- `_gt` - 大于
- `_gte` - 大于等于
- `_lt` - 小于
- `_lte` - 小于等于
- `_!` - 不等于
- `_in` - IN 查询

---

## 请求处理生命周期

### HTTP 请求处理流程

```
1. 请求到达
       │
       ▼
2. Gorilla Mux 路由匹配
       │
       ▼
3. 中间件链执行
   ┌───────────────────────────────────┐
   │ HttpContextLogMiddleware          │
   │   → 记录开始时间                    │
   │   → 包装 ResponseWriter            │
   └───────────────────────────────────┘
       │
       ▼
   ┌───────────────────────────────────┐
   │ InternalServerErrorMiddleware     │
   │   → 设置 defer recover            │
   └───────────────────────────────────┘
       │
       ▼
   ┌───────────────────────────────────┐
   │ Custom Middleware 1               │
   │   → 返回 true 继续                 │
   └───────────────────────────────────┘
       │
       ▼
4. UrlPatternHandle.Handle()
       │
       ▼
   ┌───────────────────────────────────┐
   │ 反射创建 Handler 实例               │
   └───────────────────────────────────┘
       │
       ▼
   ┌───────────────────────────────────┐
   │ InitHandler(w, r)                 │
   │   → 设置 ResponseWriter           │
   │   → 设置 Request                  │
   │   → 解析 Form                     │
   └───────────────────────────────────┘
       │
       ▼
   ┌───────────────────────────────────┐
   │ PassJson()                        │
   │   → 如果是 JSON 请求，解析 body     │
   └───────────────────────────────────┘
       │
       ▼
   ┌───────────────────────────────────┐
   │ BeforeRequest()                   │
   │   → 子类可重写的请求前钩子          │
   └───────────────────────────────────┘
       │
       ▼
   ┌───────────────────────────────────┐
   │ MethodEnum(method) 调用           │
   │   → GET → Get()                   │
   │   → POST → Post()                 │
   │   → PUT → Put()                   │
   │   → DELETE → Delete()             │
   │   → ...                           │
   └───────────────────────────────────┘
       │
       ▼
   ┌───────────────────────────────────┐
   │ TeardownRequest()                 │
   │   → 子类可重写的请求后钩子          │
   └───────────────────────────────────┘
       │
       ▼
5. 中间件 defer 执行
   ┌───────────────────────────────────┐
   │ HttpContextLogMiddleware defer    │
   │   → 计算耗时                       │
   │   → 记录日志                       │
   └───────────────────────────────────┘
       │
       ▼
6. 响应返回客户端
```

### WebSocket 连接处理流程

```
1. HTTP 请求到达
       │
       ▼
2. 路由匹配到 WSPatternHandle
       │
       ▼
3. WebSocket 升级
   ┌───────────────────────────────────┐
   │ upgrader.Upgrade(w, r, nil)       │
   │   → HTTP → WebSocket              │
   └───────────────────────────────────┘
       │
       ▼
4. Handler 初始化
   ┌───────────────────────────────────┐
   │ 反射创建 Handler 实例               │
   │ 设置 conn, request, config        │
   └───────────────────────────────────┘
       │
       ▼
5. Communicate() 执行
   ┌───────────────────────────────────┐
   │ OnConnect()                       │
   │   → 连接建立钩子                   │
   └───────────────────────────────────┘
       │
       ▼
   ┌───────────────────────────────────┐
   │ 消息循环                          │
   │   for {                           │
   │     ReadMessage()                 │
   │     OnMessage(type, data)         │
   │   }                               │
   └───────────────────────────────────┘
       │
       ▼
   ┌───────────────────────────────────┐
   │ OnDisconnect()                    │
   │   → 连接关闭钩子                   │
   └───────────────────────────────────┘
```

---

## README.md Demo 分析

### 示例代码

```go
package main

import (
    "net/http"
    "github.com/karldoenitz/Tigo/web"
)

// DemoHandler handler
type DemoHandler struct {
    web.BaseHandler  // 继承 BaseHandler
}

func (demoHandler *DemoHandler) Get() {
    demoHandler.ResponseAsText("Hello Demo!")
}

// Authorize 中间件
func Authorize(w *http.ResponseWriter, r *http.Request) bool {
    // 返回 true 继续执行，false 则停止
    return true
}

// 路由配置
var urls = []web.Pattern{
    {"/demo", DemoHandler{}, []web.Middleware{Authorize}},
}

func main() {
    application := web.Application{
        IPAddress:   "127.0.0.1",
        Port:        8888,
        UrlPatterns: urls,
    }
    application.Run()
}
```

### 代码解析

#### 1. Handler 定义

```go
type DemoHandler struct {
    web.BaseHandler
}
```

- **继承模式**：`DemoHandler` 通过嵌入 `web.BaseHandler` 继承所有基础功能
- **方法重写**：只需重写需要支持的 HTTP 方法（如 `Get()`）
- **未实现方法**：调用未重写的方法（如 `Post()`）会返回 405 Method Not Allowed

#### 2. 中间件定义

```go
func Authorize(w *http.ResponseWriter, r *http.Request) bool {
    return true
}
```

- **签名**：`func(*http.ResponseWriter, *http.Request) bool`
- **返回值**：
  - `true` - 继续执行后续中间件和 Handler
  - `false` - 中断请求处理

#### 3. 路由配置

```go
var urls = []web.Pattern{
    {"/demo", DemoHandler{}, []web.Middleware{Authorize}},
}
```

- `Url` - URL 路径
- `Handler` - Handler 实例
- `Middleware` - 中间件列表

#### 4. 应用启动

```go
application := web.Application{
    IPAddress:   "127.0.0.1",
    Port:        8888,
    UrlPatterns: urls,
}
application.Run()
```

### 请求处理流程

```
GET /demo
    │
    ▼
路由匹配: /demo → DemoHandler
    │
    ▼
中间件链:
    HttpContextLogMiddleware
        → InternalServerErrorMiddleware
            → Authorize (返回 true)
                │
                ▼
Handler 生命周期:
    InitHandler(w, r)
        → PassJson()
            → BeforeRequest()
                → Get()  ← 业务逻辑
                    → TeardownRequest()
    │
    ▼
响应: "Hello Demo!"
```

### 扩展示例

#### 带路径参数的 Handler

```go
type UserHandler struct {
    web.BaseHandler
}

func (h *UserHandler) Get() {
    userId := h.GetPathParamStr("id")
    h.ResponseAsJson(map[string]string{"user_id": userId})
}

var urls = []web.Pattern{
    {"/user/{id}", UserHandler{}, nil},
}
```

#### 带 JSON 请求体的 Handler

```go
type LoginHandler struct {
    web.BaseHandler
}

type LoginRequest struct {
    Username string `json:"username" required:"true"`
    Password string `json:"password" required:"true"`
}

func (h *LoginHandler) Post() {
    var req LoginRequest
    if err := h.CheckJsonBinding(&req); err != nil {
        h.ResponseAsJson(map[string]string{"error": err.Error()})
        return
    }
    // 业务逻辑...
    h.ResponseAsJson(map[string]string{"token": "xxx"})
}
```

#### WebSocket Handler

```go
type ChatHandler struct {
    web.WSBaseHandler
}

func (h *ChatHandler) OnConnect() {
    h.SendText("Welcome!")
}

func (h *ChatHandler) OnMessage(messageType int, data []byte) {
    // 回显消息
    h.SendText(string(data))
}

func (h *ChatHandler) OnDisconnect() {
    // 清理资源
}

var urls = []web.Pattern{
    {"/chat", ChatHandler{}, nil},
}
```

---

## 总结

Tigo web 包采用了一种简洁而强大的设计模式：

1. **Handler 继承** - 通过继承 `BaseHandler` 获得完整的请求处理能力
2. **反射驱动** - 自动根据 HTTP 方法调用对应的 Handler 方法
3. **中间件链** - 灵活的请求预处理和后处理机制
4. **插件化设计** - Session 等功能支持第三方扩展
5. **类型安全** - 通过 `ReqParams` 和 `PathParam` 提供类型转换

这种设计使得开发者可以专注于业务逻辑，而框架则处理了路由、参数解析、响应格式化等通用功能。
