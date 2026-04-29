# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Tigo is a Go web framework with a handler-based architecture. Handlers inherit from `web.BaseHandler` and implement HTTP methods (Get, Post, Put, Delete, etc.) that are automatically routed based on request method.

**Module Path:** `github.com/karldoenitz/Tigo`

## Project Structure

- `web/` - Core framework package: Application, handlers, middleware, routing, utilities
- `binding/` - JSON and form binding with validation
- `logger/` - Logging functionality
- `request/` - HTTP client for making outgoing requests
- `test_case/` - Test files
- `external_tools/tiger/` - Scaffolding CLI tool for creating Tigo projects

## Common Commands

### Building/Testing
```bash
# Run all tests
go test ./test_case/...

# Run specific test
go test ./test_case/client_test.go

# Build the main application
go build main.go
```

### Tiger CLI Tool (Scaffolding)
```bash
# Install tiger
go install github.com/karldoenitz/Tigo/external_tools/tiger@latest

# Create new Tigo project
tiger create projectName
```

## Architecture

### Handler Pattern
All handlers inherit from `web.BaseHandler` and implement HTTP methods by name:
- `Get()` - Handles GET requests
- `Post()` - Handles POST requests
- `Put()`, `Delete()`, etc. - Other HTTP methods

Unimplemented methods return 405 Method Not Allowed.

### Handler Lifecycle
1. Request is routed to handler via reflection-based method dispatch
2. `InitHandler()` is called to set up ResponseWriter and Request
3. `BeforeRequest()` hook runs (override in custom handler)
4. HTTP method handler (Get/Post/etc.) executes
5. `TeardownRequest()` hook runs (override in custom handler)

The framework automatically detects and calls the appropriate HTTP method using reflection.

### Application Lifecycle
1. Create `web.Application` with IP, Port, UrlPatterns, and optional ConfigPath
2. Call `application.Run()` (or `EndlessStart()` / `OverseerStart()` for graceful restart)
3. `InitApp()` is called internally to set up routing and load config

### Routing
Routes defined as `[]web.Pattern`:
- `Url` - URL path string
- `Handler` - Handler instance (inherits `BaseHandler`)
- `Middleware` - Optional slice of middleware functions

The framework uses `gorilla/mux` for routing under the hood, supporting path parameters via `mux.Vars()`.

Example:
```go
var urls = []web.Pattern{
    {"/demo", DemoHandler{}, []web.Middleware{Authorize}},
    {"/user/{id}", UserHandler{}, nil},
}
```

### Middleware
Custom middleware signature: `func(w *http.ResponseWriter, r *http.Request) bool`
- Return `true` to continue to next handler
- Return `false` to stop request processing

Built-in middleware (always applied):
- `HttpContextLogMiddleware` - Logs request duration and status
- `InternalServerErrorMiddleware` - Catches panics and returns 500

### WebSocket Support
Handlers with a `conn *websocket.Conn` field are automatically detected as WebSocket handlers. WebSocket connections use `web.WSBaseHandler` as base.

### File Server
Mount static files with `application.MountFileServer("/path/to/files", "/files/")`

## Configuration

Global config can be loaded from JSON or YAML file via `ConfigPath` field:
- `ip` - IP address
- `port` - Port number
- `cert` / `cert_key` - HTTPS certificate paths
- `cookie` - Cookie encryption key
- `template` - Template file directory path
- `log` - Log level configuration

## Handler Methods

`BaseHandler` provides common methods:
- **Response**: `ResponseAsJson()`, `ResponseAsText()`, `ResponseAsHtml()`, `Render()` (templates)
- **Parameters**: `GetParameter()`, `GetPathParam()`, `GetJsonValue()`, `CheckParamBinding()`
- **Cookies**: `SetCookie()`, `GetCookie()`, `SetSecureCookie()`, `ClearCookie()`
- **Session**: `SetSession()`, `GetSession()`, `DelSession()` (requires `application.StartSession()`)
- **Headers**: `GetHeader()`, `SetHeader()`
- **Redirects**: `Redirect()`, `Move()`, `RedirectPermanently()`
- **Context**: `SetCtxVal()`, `GetCtxVal()` for request-scoped values

## Cookie Encryption

Tigo provides built-in cookie encryption via the `web.Cookie` struct:
- Use `SetSecureCookie()` / `GetSecureCookie()` for encrypted cookies
- Set encryption key via `GlobalConfig.Cookie` or pass as parameter
- Supports advanced cookie options with `SetAdvancedCookie()` (path, domain, expires, secure, httpOnly)

## Custom Validation

Structs used for binding can implement a `Check() error` method for custom validation logic that runs after tag-based validation.

## Binding/Validation

Use `binding` package for structured request validation:
- `json:"field_name"` tags for JSON
- `form:"field_name"` tags for form data
- Call `baseHandler.CheckJsonBinding(obj)` or `CheckFormBinding(obj)` in handlers

**Validation Tags:**
- `required:"true"` - marks field as required
- `default:"value"` - sets default value if field is empty
- `regex:"pattern"` - validates field against regex pattern

Example:
```go
type UserRequest struct {
    Name  string `json:"name" required:"true"`
    Email string `json:"email" required:"true" regex:"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+$"`
    Age   int    `json:"age" default:"18"`
}
```

## Session Support

Tigo uses a pluggable session architecture. Third-party session implementations must implement the `web.SessionInterface`:

```go
type SessionInterface interface {
    NewSessionManager() SessionManager
}
```

To enable sessions, call `application.StartSession()` with a session manager before running the app. For a ready-to-use Redis/MySQL session plugin, see [tission](https://github.com/karldoenitz/tission).

## Third-Party Tools

- **tiger** - Scaffolding CLI tool for creating Tigo projects (`go install github.com/karldoenitz/Tigo/external_tools/tiger@latest`)
- **tission** - Session plugin for Redis/MySQL backend storage

## Test Files

Tests are located in `test_case/` directory. Run with standard Go test commands.