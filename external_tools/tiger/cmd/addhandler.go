package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/cli"
	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/fileutil"
	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/template"
	"github.com/spf13/cobra"
)

var (
	handlerPath    string
	handlerRoute   string
	handlerMethods []string
)

var addHandlerCmd = &cobra.Command{
	Use:   "addhandler <handler_name>",
	Short: "添加一个新的 Handler",
	Long: `在当前项目中添加一个新的 Handler。
Handler 将被添加到 handler 目录，并自动注册路由。`,
	Args: cobra.ExactArgs(1),
	RunE: runAddHandler,
}

func init() {
	rootCmd.AddCommand(addHandlerCmd)

	addHandlerCmd.Flags().StringVarP(&handlerPath, "path", "p", "handler", "Handler 目录路径")
	addHandlerCmd.Flags().StringVar(&handlerRoute, "route", "", "路由路径 (如: /user)")
	addHandlerCmd.Flags().StringSliceVar(&handlerMethods, "methods", []string{"GET"}, "HTTP 方法 (GET, POST, PUT, DELETE)")
}

func runAddHandler(cmd *cobra.Command, args []string) error {
	handlerName := args[0]

	// 验证 handler 名称
	if err := cli.ValidateProjectName(handlerName); err != nil {
		return fmt.Errorf("Handler 名称无效: %w", err)
	}

	// 检查是否在 Tigo 项目中
	workDir, err := fileutil.GetWorkingDir()
	if err != nil {
		return err
	}

	if !isTigoProject(workDir) {
		return fmt.Errorf("当前目录不是 Tigo 项目，请在项目根目录中运行此命令")
	}

	// 初始化模板引擎
	engine := template.NewEngine()

	// 注册 handler 模板
	handlerTemplate := getHandlerTemplate(handlerName)
	if err := engine.Register("handler", handlerTemplate); err != nil {
		return err
	}

	// 创建 handler 目录
	handlerDir := filepath.Join(workDir, handlerPath)
	if err := fileutil.EnsureDir(handlerDir); err != nil {
		return err
	}

	// 渲染 handler 文件
	data := template.TemplateData{
		"HandlerName": handlerName,
	}
	handlerFile := filepath.Join(handlerDir, toFileName(handlerName))
	if err := engine.RenderToFile("handler", data, handlerFile); err != nil {
		return err
	}

	fmt.Printf("Handler %s 创建成功: %s\n", handlerName, handlerFile)

	// 添加路由到 main.go
	if err := addRouteToMain(workDir, handlerName, handlerRoute); err != nil {
		return err
	}

	return nil
}

// getHandlerTemplate 获取 handler 模板
func getHandlerTemplate(name string) string {
	return `// {{ .HandlerName }} Handler
package handler

import (
	"github.com/karldoenitz/Tigo/web"
)

type {{ .HandlerName }} struct {
	web.BaseHandler
}

// Get 处理 GET 请求
func (h *{{ .HandlerName }}) Get() {
	h.ResponseAsText("{{ .HandlerName }} - GET")
}

// Post 处理 POST 请求
func (h *{{ .HandlerName }}) Post() {
	h.ResponseAsText("{{ .HandlerName }} - POST")
}

// Put 处理 PUT 请求
func (h *{{ .HandlerName }}) Put() {
	h.ResponseAsText("{{ .HandlerName }} - PUT")
}

// Delete 处理 DELETE 请求
func (h *{{ .HandlerName }}) Delete() {
	h.ResponseAsText("{{ .HandlerName }} - DELETE")
}
`
}

// isTigoProject 检查是否是 Tigo 项目
func isTigoProject(dir string) bool {
	// 检查是否有 main.go 和 go.mod
	mainExists := fileutil.Exists(filepath.Join(dir, "main.go"))
	goModExists := fileutil.Exists(filepath.Join(dir, "go.mod"))
	return mainExists && goModExists
}

// toFileName 将 handler 名转换为文件名
func toFileName(name string) string {
	return cli.SanitizeString(name) + ".go"
}

// addRouteToMain 添加路由到 main.go
// TODO: 使用 go/ast 安全地修改代码
func addRouteToMain(projectDir, handlerName, route string) error {
	mainPath := filepath.Join(projectDir, "main.go")

	// 生成路由行
	routePath := route
	if routePath == "" {
		// 移除 "Handler" 后缀
		nameWithoutSuffix := handlerName
		if strings.HasSuffix(handlerName, "Handler") {
			nameWithoutSuffix = handlerName[:len(handlerName)-7]
		}
		routePath = "/" + toFileName(nameWithoutSuffix)
	}

	importLine := `"` + getProjectName(projectDir) + `/handler"` + "\n"
	routeLine := `{"` + routePath + `", handler.` + handlerName + `{}, nil},` + "\n"

	// 简单实现：查找并插入
	lines, err := fileutil.ReadLines(mainPath)
	if err != nil {
		return err
	}

	var newLines []string
	importAdded := false
	routeAdded := false

	for _, line := range lines {
		// 添加 import
		if !importAdded && strings.Contains(line, "github.com/karldoenitz/Tigo/web") {
			newLines = append(newLines, importLine)
			importAdded = true
		}
		newLines = append(newLines, line)

		// 添加路由
		if !routeAdded && strings.Contains(line, "var urls = []web.Pattern{") {
			newLines = append(newLines, routeLine)
			routeAdded = true
		}
	}

	if err := fileutil.WriteLines(mainPath, newLines); err != nil {
		return err
	}

	fmt.Printf("路由已添加: %s -> handler.%s\n", routePath, handlerName)
	return nil
}

// getProjectName 从目录路径获取项目名
func getProjectName(dir string) string {
	return filepath.Base(dir)
}
