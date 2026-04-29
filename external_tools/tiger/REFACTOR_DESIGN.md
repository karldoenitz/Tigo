# Tiger 脚手架工具重构设计文档

## 1. 现状分析

### 1.1 现有功能
- `create <project_name>` - 创建 Tigo 项目
- `addhandler <handler_name>` - 添加 handler
- `conf <filename>` - 添加配置文件
- `logger` - 添加日志
- `mod` - 执行 go mod
- `version` - 显示版本
- `help` - 帮助信息

### 1.2 存在的问题

| 问题            | 描述                          | 影响        |
|---------------|-----------------------------|-----------|
| 代码组织混乱        | 所有逻辑都在一个 main.go 文件中（362 行） | 难以维护和扩展   |
| 硬编码模板         | 模板代码直接作为字符串常量               | 无法自定义模板   |
| 命令解析简陋        | 使用简单的 switch case           | 功能受限，体验差  |
| 错误处理差         | 大量使用 panic                  | 遇到错误直接崩溃  |
| 不可扩展          | 添加新功能需要修改多处代码               | 开发效率低     |
| 无配置文件支持       | 无法自定义模板或配置                  | 灵活性差      |
| main.go 修改不可靠 | 使用字符串分割来修改代码                | 容易破坏原有代码  |
| 缺少交互式输入       | 只能通过命令行参数                   | 用户体验差     |
| 无模板版本管理       | 模板无法更新                      | 无法跟进框架新特性 |
| 缺少测试          | 几乎没有单元测试                    | 代码质量无保障   |
| 代码重复          | 很多相似逻辑重复                    | 违反 DRY 原则 |

## 2. 重构目标

1. **模块化** - 清晰的代码结构和职责分离
2. **可扩展** - 插件化命令系统，易于添加新功能
3. **用户友好** - 交互式体验，清晰的错误提示
4. **可配置** - 支持自定义模板和配置
5. **健壮性** - 完善的错误处理和测试覆盖
6. **现代 CLI 标准** - 使用成熟的 CLI 框架

## 3. 架构设计

### 3.1 目录结构

```
external_tools/tiger/
├── cmd/
│   ├── root.go              # 根命令
│   ├── create.go            # create 子命令
│   ├── addhandler.go        # addhandler 子命令
│   ├── conf.go              # conf 子命令
│   ├── logger.go            # logger 子命令
│   ├── mod.go               # mod 子命令
│   ├── template.go          # template 子命令
│   └── version.go           # version 子命令
├── pkg/
│   ├── cli/                 # CLI 核心逻辑
│   │   ├── command.go       # 命令接口定义
│   │   ├── runner.go        # 命令执行器
│   │   └── prompt.go        # 交互式输入
│   ├── template/            # 模板管理
│   │   ├── engine.go        # 模板引擎
│   │   ├── registry.go      # 模板注册表
│   │   └── parser.go        # 模板解析器
│   ├── project/             # 项目操作
│   │   ├── creator.go       # 项目创建
│   │   └── modifier.go      # 项目修改
│   ├── fileutil/            # 文件工具
│   │   ├── writer.go        # 文件写入
│   │   ├── reader.go        # 文件读取
│   │   └── parser.go        # 代码解析/修改
│   ├── config/              # 配置管理
│   │   ├── loader.go        # 配置加载
│   │   └── validator.go     # 配置验证
│   └── git/                 # Git 操作
│       └── git.go           # Git 封装
├── templates/               # 模板文件目录
│   ├── project/
│   │   ├── main.go.tmpl     # main.go 模板
│   │   ├── handler.go.tmpl # handler 模板
│   │   └── config.json.tmpl # 配置模板
│   ├── api/                 # API 模板
│   │   ├── restful.go.tmpl
│   │   └── graphql.go.tmpl
│   └── middleware/          # 中间件模板
│       └── auth.go.tmpl
├── config/
│   └── tiger.yaml           # Tiger 配置文件
├── main.go                  # 入口文件
└── go.mod                   # 依赖管理
```

### 3.2 核心组件

#### 3.2.1 命令系统 (cmd/)

采用 Cobra CLI 框架，提供：
- 子命令嵌套支持
- 标志参数解析
- 自动生成帮助文档
- Shell 自动补全

#### 3.2.2 模板引擎 (pkg/template/)

```go
type TemplateEngine interface {
    // 渲染模板
    Render(name string, data interface{}) ([]byte, error)
    // 注册模板
    Register(name string, content string) error
    // 加载模板目录
    LoadTemplates(dir string) error
    // 列出可用模板
    List() []string
}
```

支持：
- Go template 语法
- 嵌套模板
- 自定义函数
- 多种格式输出

#### 3.2.3 项目创建器 (pkg/project/)

```go
type ProjectCreator interface {
    // 创建项目
    Create(name string, opts CreateOptions) error
    // 验证项目名称
    ValidateName(name string) error
    // 检查目录是否存在
    CheckDirExists(path string) bool
}

type CreateOptions struct {
    Template      string
    IncludeTests  bool
    IncludeConfig bool
    GitInit       bool
}
```

#### 3.2.4 代码修改器 (pkg/fileutil/parser/)

使用 Go 的 `go/parser` 和 `go/ast` 来安全地修改代码：

```go
type CodeModifier interface {
    // 添加导入
    AddImport(filePath string, imp string) error
    // 添加路由
    AddRoute(filePath string, route Route) error
    // 添加处理器
    AddHandler(filePath string, handler string) error
}

type Route struct {
    Path    string
    Handler string
    Methods []string
}
```

#### 3.2.5 交互式输入 (pkg/cli/prompt/)

使用 survey 库提供交互式体验：

```go
type Prompter interface {
    // 询问字符串输入
    AskString(question, default string) (string, error)
    // 询问选项选择
    AskSelect(question string, options []string) (string, error)
    // 确认操作
    Confirm(question string, default bool) (bool, error)
    // 多选
    AskMultiSelect(question string, options []string) ([]string, error)
}
```

## 4. 功能设计

### 4.1 命令增强

#### create - 项目创建

```bash
# 基础用法
tiger create myapp

# 使用模板
tiger create myapp --template restful

# 交互式创建
tiger create myapp --interactive

# 指定选项
tiger create myapp \n  --template restful \n  --port 8080 \n  --enable-logger \n  --enable-session \n  --git-init
```

支持的模板：
- `basic` - 基础项目
- `restful` - RESTful API 项目
- `websocket` - WebSocket 项目
- `graphql` - GraphQL 项目
- `full` - 全功能项目

#### addhandler - 添加 Handler

```bash
# 基础用法
tiger addhandler UserHandler

# 指定路由
tiger addhandler UserHandler --route /user

# 指定 HTTP 方法
tiger addhandler UserHandler --methods GET,POST,PUT,DELETE

# 添加到指定目录
tiger addhandler UserHandler --path internal/handlers

# 添加到中间件后
tiger addhandler UserHandler --after AuthMiddleware
```

#### template - 模板管理（新命令）

```bash
# 列出可用模板
tiger template list

# 安装模板
tiger template install github.com/user/tigo-templates

# 更新模板
tiger template update

# 创建自定义模板
tiger template create mytemplate

# 删除模板
tiger template remove mytemplate
```

#### conf - 配置管理

```bash
# 生成配置文件
tiger conf generate --format json
tiger conf generate --format yaml

# 验证配置
tiger conf validate config.yaml

# 查看配置
tiger conf show config.yaml
```

### 4.2 新增功能

#### 插件系统

```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    Execute(ctx Context) error
}
```

支持的插件：
- 数据库迁移插件
- 中间件生成插件
- 测试生成插件
- 文档生成插件

#### 项目健康检查

```bash
tiger doctor
```

检查项：
- Go 版本兼容性
- 依赖完整性
- 代码规范
- 安全配置

#### 配置向导

```bash
tiger wizard
```

引导式配置项目，自动生成配置文件。

## 5. 技术选型

| 组件 | 技术选择 | 理由 |
|------|---------|------|
| CLI 框架 | Cobra | 成熟、广泛使用、功能完善 |
| 交互式输入 | survey | 友好的交互体验 |
| 代码解析 | go/ast | 官方库、准确可靠 |
| 配置文件 | viper | 支持多种格式 |
| 日志 | zap | 高性能、结构化 |
| 测试 | testify | 丰富的断言和 mock |
| 颜色输出 | lipgloss | 现代化的终端输出 |

## 6. 向后兼容性

保持现有命令的行为不变，但增加更多选项：

```bash
# 旧用法仍然有效
tiger create myapp
tiger addhandler UserHandler
tiger conf config.json

# 新增的选项不影响旧用法
tiger create myapp --template restful
```

## 7. 错误处理策略

```go
type TigerError struct {
    Code    int
    Message string
    Details string
    Cause   error
}

func (e *TigerError) Error() string {
    // 格式化错误信息
}

// 错误代码定义
const (
    ErrProjectExists = 1001
    ErrInvalidName    = 1002
    ErrTemplateNotFound = 1003
    ErrFileWrite      = 1004
    // ...
)
```

提供清晰的错误信息和建议的解决方案。

## 8. 测试策略

```go
// 测试示例
func TestCreateProject(t *testing.T) {
    // 使用临时目录
    tmpDir := t.TempDir()

    creator := NewProjectCreator()
    err := creator.Create("testapp", CreateOptions{
        BaseDir: tmpDir,
    })

    assert.NoError(t, err)
    assert.FileExists(t, tmpDir + "/testapp/main.go")
}
```

覆盖：
- 单元测试
- 集成测试
- 端到端测试
- 模板渲染测试

## 9. 性能考虑

- 模板预编译
- 文件操作使用缓冲
- 并行处理不相关的操作
- 缓存模板解析结果

## 10. 文档

- README.md - 基本使用
- docs/commands.md - 命令详解
- docs/templates.md - 模板开发
- docs/plugins.md - 插件开发
- docs/migration.md - 迁移指南

## 11. 发布计划

### Phase 1: 核心重构
- [ ] 重构现有代码结构
- [ ] 迁移到 Cobra
- [ ] 实现基础模板引擎

### Phase 2: 功能增强
- [ ] 添加交互式输入
- [ ] 实现代码安全修改
- [ ] 添加模板管理命令

### Phase 3: 插件系统
- [ ] 设计插件接口
- [ ] 实现插件加载
- [ ] 开发官方插件

### Phase 4: 完善生态
- [ ] 完善文档
- [ ] 添加测试
- [ ] 发布稳定版本

## 12. 配置示例

### tiger.yaml

```yaml
# Tiger 配置文件
version: "2.0"

# 默认模板
default_template: "basic"

# 模板仓库
templates:
  - name: "basic"
    path: "./templates/project/basic"
  - name: "restful"
    path: "./templates/project/restful"
  - name: "websocket"
    path: "./templates/project/websocket"

# 模板仓库 URL
template_repositories:
  - "https://github.com/karldoenitz/tigo-templates"

# 项目默认配置
project_defaults:
  port: 8080
  ip: "0.0.0.0"
  enable_logger: true
  enable_session: false

# 编辑器配置
editor:
  format: "gofmt"
  imports: "goimports"

# Git 配置
git:
  auto_init: false
  default_branch: "main"
```

## 13. 总结

通过本次重构，Tiger 脚手架工具将实现：

1. ✅ 清晰的代码结构
2. ✅ 良好的扩展性
3. ✅ 友好的用户体验
4. ✅ 可靠的代码修改
5. ✅ 完善的测试覆盖
6. ✅ 丰富的模板生态

这将使 Tiger 成为 Tigo 框架开发者不可或缺的生产力工具。