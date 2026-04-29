package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/cli"
	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/template"
)

// CreateOptions 项目创建选项
type CreateOptions struct {
	Name          string
	Template      string
	Port          int
	IP            string
	EnableLogger  bool
	EnableSession bool
	IncludeTests  bool
	GitInit       bool
	Interactive   bool
}

// Creator 项目创建器
type Creator struct {
	prompter cli.Prompter
	engine   *template.Engine
}

// NewCreator 创建新的项目创建器
func NewCreator(engine *template.Engine, prompter cli.Prompter) *Creator {
	return &Creator{
		prompter: prompter,
		engine:   engine,
	}
}

// Create 创建项目
func (c *Creator) Create(opts CreateOptions) error {
	// 验证项目名称
	if err := cli.ValidateProjectName(opts.Name); err != nil {
		return err
	}

	// 检查目录是否已存在
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}
	projectPath := filepath.Join(workDir, opts.Name)

	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("目录 %s 已存在", projectPath)
	}

	// 如果是交互模式，询问用户选项
	if opts.Interactive {
		var err error
		opts, err = c.interactiveCreate(opts.Name)
		if err != nil {
			return err
		}
	}

	// 设置默认值
	if opts.Template == "" {
		opts.Template = "basic"
	}
	if opts.Port == 0 {
		opts.Port = 8080
	}
	if opts.IP == "" {
		opts.IP = "0.0.0.0"
	}

	// 创建项目目录
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("创建项目目录失败: %w", err)
	}

	// 渲染模板
	if err := c.renderTemplates(opts, projectPath); err != nil {
		return fmt.Errorf("渲染模板失败: %w", err)
	}

	// 初始化 Git
	if opts.GitInit {
		if err := c.initGit(projectPath); err != nil {
			return fmt.Errorf("初始化 Git 失败: %w", err)
		}
	}

	fmt.Printf("\n项目 %s 创建成功!\n", opts.Name)
	fmt.Printf("进入项目目录: cd %s\n", opts.Name)
	fmt.Println("运行以下命令启动项目:")
	fmt.Printf("  go mod init %s\n", opts.Name)
	fmt.Printf("  go mod tidy\n")
	fmt.Println("  go run main.go")

	return nil
}

// interactiveCreate 交互式创建项目
func (c *Creator) interactiveCreate(name string) (CreateOptions, error) {
	opts := CreateOptions{Name: name}

	// 选择模板
	templateName, err := c.prompter.AskSelect(
		"选择项目模板",
		[]string{"basic", "restful", "websocket", "graphql", "full"},
	)
	if err != nil {
		return opts, err
	}
	opts.Template = templateName

	// 询问端口号
	portStr, err := c.prompter.AskString("端口号 (默认: 8080)", "8080")
	if err != nil {
		return opts, err
	}
	fmt.Sscanf(portStr, "%d", &opts.Port)

	// 询问 IP 地址
	opts.IP, err = c.prompter.AskString("IP 地址 (默认: 0.0.0.0)", "0.0.0.0")
	if err != nil {
		return opts, err
	}

	// 询问是否启用日志
	opts.EnableLogger, err = c.prompter.Confirm("启用日志", true)
	if err != nil {
		return opts, err
	}

	// 询问是否启用 Session
	opts.EnableSession, err = c.prompter.Confirm("启用 Session", false)
	if err != nil {
		return opts, err
	}

	// 询问是否初始化 Git
	opts.GitInit, err = c.prompter.Confirm("初始化 Git", false)
	if err != nil {
		return opts, err
	}

	return opts, nil
}

// renderTemplates 渲染模板文件
func (c *Creator) renderTemplates(opts CreateOptions, projectPath string) error {
	data := template.TemplateData{
		"ProjectName":   opts.Name,
		"PackageName":   toPackageName(opts.Name),
		"Port":          opts.Port,
		"IP":            opts.IP,
		"EnableLogger":  opts.EnableLogger,
		"EnableSession": opts.EnableSession,
		"IncludeTests":  opts.IncludeTests,
	}

	// 渲染 main.go
	mainPath := filepath.Join(projectPath, "main.go")
	if err := c.engine.RenderToFile(opts.Template+"/main", data, mainPath); err != nil {
		return err
	}

	// 创建 handler 目录
	handlerDir := filepath.Join(projectPath, "handler")
	if err := os.MkdirAll(handlerDir, 0755); err != nil {
		return err
	}

	// 渲染 ping handler
	pingPath := filepath.Join(handlerDir, "pinghandler.go")
	if err := c.engine.RenderToFile(opts.Template+"/handler", data, pingPath); err != nil {
		return err
	}

	// 如果需要测试文件
	if opts.IncludeTests {
		testPath := filepath.Join(handlerDir, "pinghandler_test.go")
		if err := c.engine.RenderToFile(opts.Template+"/handler_test", data, testPath); err != nil {
			return err
		}
	}

	// 如果启用日志，创建日志配置
	if opts.EnableLogger {
		logDir := filepath.Join(projectPath, "log")
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// initGit 初始化 Git 仓库
func (c *Creator) initGit(projectPath string) error {
	// 检查 git 命令是否可用
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git 命令未找到，请先安装 git")
	}

	// TODO: 实现实际的 git init 命令
	// 这里需要使用 exec.Command 来运行 git init

	return nil
}

// toPackageName 将项目名转换为包名
func toPackageName(name string) string {
	// 移除非法字符
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	name = reg.ReplaceAllString(name, "_")

	// 确保不以数字开头
	if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
		name = "_" + name
	}

	// 转换为小写
	name = strings.ToLower(name)

	return name
}
