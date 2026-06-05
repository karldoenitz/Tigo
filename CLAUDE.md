# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 提供项目代码指导。

## 概述

Tigo 是一个基于 Handler 架构的 Go Web 框架。Handler 继承自 `web.BaseHandler`，通过实现 HTTP 方法（Get、Post、Put、Delete 等）来自动路由请求。

**模块路径:** `github.com/karldoenitz/Tigo`

## 项目结构

- `web/` - 核心框架包：Application、Handler、中间件、路由、工具函数
- `binding/` - JSON 和表单数据绑定及验证
- `logger/` - 日志功能
- `request/` - HTTP 客户端（用于发起外部请求）
- `test_case/` - 测试文件
- `external_tools/tiger/` - 脚手架 CLI 工具，用于创建 Tigo 项目

## 常用命令

### 构建与测试
```bash
# 运行所有测试
go test ./test_case/...

# 运行特定测试文件
go test ./test_case/client_test.go

# 构建主程序
go build main.go
```

### Tiger CLI 工具（脚手架）
```bash
# 安装 tiger
go install github.com/karldoenitz/Tigo/external_tools/tiger@latest

# 创建新的 Tigo 项目
tiger create projectName
```

## 架构设计

### Handler 模式
所有 Handler 继承 `web.BaseHandler`，通过方法名实现对应的 HTTP 方法：
- `Get()` - 处理 GET 请求
- `Post()` - 处理 POST 请求
- `Put()`、`Delete()` 等 - 处理其他 HTTP 方法

未实现的方法自动返回 405 Method Not Allowed。

### Handler 生命周期
1. 请求通过反射机制路由到对应的 Handler
2. 调用 `InitHandler()` 初始化 ResponseWriter 和 Request
3. 执行 `BeforeRequest()` 前置钩子（可在自定义 Handler 中重写）
4. 执行对应的 HTTP 方法处理函数（Get/Post 等）
5. 执行 `TeardownRequest()` 后置钩子（可在自定义 Handler 中重写）

框架通过反射自动检测并调用对应的 HTTP 方法。

### Application 生命周期
1. 创建 `web.Application`，传入 IP、Port、UrlPatterns 和可选的 ConfigPath
2. 调用 `application.Run()`（或 `EndlessStart()` / `OverseerStart()` 实现优雅重启）
3. 内部调用 `InitApp()` 初始化路由并加载配置

### 路由
路由定义为 `[]web.Pattern`：
- `Url` - URL 路径字符串
- `Handler` - Handler 实例（继承 `BaseHandler`）
- `Middleware` - 可选的中间件函数切片

底层使用 `gorilla/mux` 进行路由匹配，支持通过 `mux.Vars()` 获取路径参数。

示例：
```go
var urls = []web.Pattern{
    {"/demo", DemoHandler{}, []web.Middleware{Authorize}},
    {"/user/{id}", UserHandler{}, nil},
}
```

### 中间件
自定义中间件签名：`func(w *http.ResponseWriter, r *http.Request) bool`
- 返回 `true` 继续执行后续处理
- 返回 `false` 中止请求处理

内置中间件（始终生效）：
- `HttpContextLogMiddleware` - 记录请求耗时和状态
- `InternalServerErrorMiddleware` - 捕获 panic 并返回 500 错误

### WebSocket 支持
包含 `conn *websocket.Conn` 字段的 Handler 会被自动识别为 WebSocket Handler。WebSocket 连接使用 `web.WSBaseHandler` 作为基类。

### 静态文件服务
使用 `application.MountFileServer("/path/to/files", "/files/")` 挂载静态文件。

## 配置

通过 `ConfigPath` 字段可从 JSON 或 YAML 文件加载全局配置：
- `ip` - IP 地址
- `port` - 端口号
- `cert` / `cert_key` - HTTPS 证书路径
- `cookie` - Cookie 加密密钥
- `template` - 模板文件目录路径
- `log` - 日志级别配置

## Handler 方法

`BaseHandler` 提供的常用方法：
- **响应**: `ResponseAsJson()`、`ResponseAsText()`、`ResponseAsHtml()`、`Render()`（模板渲染）
- **参数获取**: `GetParameter()`、`GetPathParam()`、`GetJsonValue()`、`CheckParamBinding()`
- **Cookie**: `SetCookie()`、`GetCookie()`、`SetSecureCookie()`、`ClearCookie()`
- **Session**: `SetSession()`、`GetSession()`、`DelSession()`（需先调用 `application.StartSession()`）
- **请求头**: `GetHeader()`、`SetHeader()`
- **重定向**: `Redirect()`、`Move()`、`RedirectPermanently()`
- **上下文**: `SetCtxVal()`、`GetCtxVal()` 用于请求级别的值传递

## Cookie 加密

Tigo 通过 `web.Cookie` 结构体提供内置的 Cookie 加密功能：
- 使用 `SetSecureCookie()` / `GetSecureCookie()` 操作加密 Cookie
- 通过 `GlobalConfig.Cookie` 或参数传入设置加密密钥
- 支持通过 `SetAdvancedCookie()` 设置高级选项（path、domain、expires、secure、httpOnly）

## 数据绑定与验证

使用 `binding` 包进行结构化的请求数据验证：
- `json:"field_name"` 标签用于 JSON 数据
- `form:"field_name"` 标签用于表单数据
- 在 Handler 中调用 `baseHandler.CheckJsonBinding(obj)` 或 `CheckFormBinding(obj)`

**验证标签：**
- `required:"true"` - 标记字段为必填
- `default:"value"` - 字段为空时设置默认值
- `regex:"pattern"` - 根据正则表达式验证字段

示例：
```go
type UserRequest struct {
    Name  string `json:"name" required:"true"`
    Email string `json:"email" required:"true" regex:"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+$"`
    Age   int    `json:"age" default:"18"`
}
```

### 自定义验证

用于绑定的结构体可以实现 `Check() error` 方法，在标签验证之后执行自定义验证逻辑。

## Session 支持

Tigo 采用可插拔的 Session 架构。第三方 Session 实现需实现 `web.SessionInterface` 接口：

```go
type SessionInterface interface {
    NewSessionManager() SessionManager
}
```

启用 Session 需在应用启动前调用 `application.StartSession()` 并传入 Session 管理器。可使用现成的 Redis/MySQL Session 插件 [tission](https://github.com/karldoenitz/tission)。

## 相关工具

- **tiger** - 脚手架 CLI 工具，用于快速创建 Tigo 项目（`go install github.com/karldoenitz/Tigo/external_tools/tiger@latest`）
- **tission** - Session 插件，支持 Redis/MySQL 后端存储

## 测试

测试文件位于 `test_case/` 目录，使用标准 Go 测试命令运行。
