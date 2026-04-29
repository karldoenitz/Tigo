# Tiger 脚手架工具 v2.0

Tigo Web 框架的官方脚手架工具，用于快速创建项目、添加 Handler、生成配置等。

## 新版特性

- 模块化架构，易于维护和扩展
- 支持多种项目模板
- 交互式项目创建
- 使用 Cobra CLI 框架，体验更好
- 模板引擎，支持自定义模板
- 完善的错误处理
- 支持通过配置文件自定义默认选项

## 安装

```bash
# 从本地安装
go install github.com/karldoenitz/Tigo/external_tools/tiger@latest

# 或直接编译
cd external_tools/tiger
go build -o tiger .
```

## 使用方法

### 创建项目

```bash
# 基础项目
tiger create myapp

# 使用模板
tiger create myapp --template restful
tiger create myapp --template websocket

# 交互式创建
tiger create myapp --interactive

# 指定选项
tiger create myapp --port 8080 --ip 0.0.0.0 --enable-logger --enable-session --git-init
```

### 添加 Handler

```bash
tiger addhandler UserHandler
tiger addhandler UserHandler --route /user
tiger addhandler UserHandler --methods GET,POST,PUT,DELETE --path internal/handlers
```

### 生成配置文件

```bash
tiger conf config.json
tiger conf config.yaml
```

### 执行 Go Mod

```bash
tiger mod
```

### 查看版本

```bash
tiger version
```

### 获取帮助

```bash
tiger --help
tiger create --help
```

## 项目模板

- `basic` - 基础 Tigo 项目
- `restful` - RESTful API 项目
- `websocket` - WebSocket 项目
- `graphql` - GraphQL 项目（待实现）
- `full` - 全功能项目（待实现）

## 项目结构

```
external_tools/tiger/
├── cmd/                  # CLI 命令
│   ├── root.go          # 根命令
│   ├── create.go        # create 子命令
│   ├── addhandler.go    # addhandler 子命令
│   ├── conf.go          # conf 子命令
│   ├── mod.go           # mod 子命令
│   └── version.go       # version 子命令
├── pkg/                  # 核心包
│   ├── cli/            # CLI 工具
│   ├── template/        # 模板引擎
│   ├── project/         # 项目创建/修改
│   ├── fileutil/        # 文件操作工具
│   └── config/         # 配置管理
├── templates/            # 模板文件
│   ├── project/         # 项目模板
│   └── middleware/      # 中间件模板
├── main.go              # 入口文件
├── go.mod               # Go 模块定义
└── tiger.yaml            # Tiger 配置文件
```

## 配置文件

在 `tiger.yaml` 中配置默认选项：

```yaml
version: "2.0"
default_template: "basic"
templates:
  - name: "basic"
    path: "./project/basic"
  - name: "restful"
    path: "./project/restful"
  - name: "websocket"
    path: "./project/websocket"
project_defaults:
  port: 8080
  ip: "0.0.0.0"
  enable_logger: true
  enable_session: false
```

## 命令参考

| 命令 | 描述 |
|-------|------|
| `create <name>` | 创建新项目 |
| `addhandler <name>` | 添加 Handler |
| `conf <filename>` | 生成配置文件 |
| `mod` | 执行 go mod 命令 |
| `version` | 显示版本信息 |
| `help` | 显示帮助信息 |

## 与 v1.x 的兼容性

v2.0 保持了向后兼容性，原有的命令用法仍然有效：

```bash
# 旧命令仍然可用
tiger create myproject
tiger addhandler MyHandler
tiger conf config.json
tiger mod
tiger version
```

## 开发

```bash
cd external_tools/tiger
go run main.go --help
go build -o tiger .
```

## 更多信息

- [Tigo 官方文档](https://github.com/karldoenitz/Tigo)
- [问题反馈](https://github.com/karldoenitz/Tigo/issues)
- [重构设计文档](REFACTOR_DESIGN.md)