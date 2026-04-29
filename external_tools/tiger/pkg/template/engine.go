package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TemplateData 模板数据
type TemplateData map[string]interface{}

// Engine 模板引擎
type Engine struct {
	templates map[string]*template.Template
	funcMap   template.FuncMap
}

// NewEngine 创建新的模板引擎
func NewEngine() *Engine {
	return &Engine{
		templates: make(map[string]*template.Template),
		funcMap:   defaultFuncMap(),
	}
}

// defaultFuncMap 返回默认的模板函数
func defaultFuncMap() template.FuncMap {
	c := cases.Title(language.English)
	return template.FuncMap{
		"toLower":   strings.ToLower,
		"toUpper":   strings.ToUpper,
		"title":     c.String,
		"trim":      strings.TrimSpace,
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"join":      strings.Join,
		"replace":   strings.ReplaceAll,
		"split":     strings.Split,
	}
}

// AddFunc 添加自定义函数
func (e *Engine) AddFunc(name string, fn interface{}) {
	e.funcMap[name] = fn
}

// Register 注册模板
func (e *Engine) Register(name string, content string) error {
	tmpl, err := template.New(name).Funcs(e.funcMap).Parse(content)
	if err != nil {
		return fmt.Errorf("解析模板 %s 失败: %w", name, err)
	}
	e.templates[name] = tmpl
	return nil
}

// LoadFromFile 从文件加载模板
func (e *Engine) LoadFromFile(name, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %w", err)
	}
	return e.Register(name, string(content))
}

// LoadFromDir 从目录加载所有模板
func (e *Engine) LoadFromDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// 计算相对路径作为模板名
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		// 移除 .tmpl 扩展名
		name := strings.TrimSuffix(relPath, ".tmpl")
		name = strings.ReplaceAll(name, string(filepath.Separator), "/")

		return e.LoadFromFile(name, path)
	})
}

// Render 渲染模板
func (e *Engine) Render(name string, data TemplateData) ([]byte, error) {
	tmpl, ok := e.templates[name]
	if !ok {
		return nil, fmt.Errorf("模板 %s 不存在", name)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("渲染模板 %s 失败: %w", name, err)
	}

	return buf.Bytes(), nil
}

// RenderToFile 渲染模板到文件
func (e *Engine) RenderToFile(name string, data TemplateData, outputPath string) error {
	content, err := e.Render(name, data)
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// List 列出所有已注册的模板
func (e *Engine) List() []string {
	names := make([]string, 0, len(e.templates))
	for name := range e.templates {
		names = append(names, name)
	}
	return names
}

// Exists 检查模板是否存在
func (e *Engine) Exists(name string) bool {
	_, ok := e.templates[name]
	return ok
}
