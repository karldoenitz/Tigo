package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/fileutil"
	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/template"
	"github.com/karldoenitz/Tigo/web"
	"github.com/spf13/cobra"
)

var confFormat string

var confCmd = &cobra.Command{
	Use:   "conf <filename>",
	Short: "生成配置文件",
	Long:  `生成 Tigo 项目配置文件 (JSON 或 YAML 格式)。`,
	Args:  cobra.ExactArgs(1),
	RunE:  runConf,
}

func init() {
	rootCmd.AddCommand(confCmd)
	confCmd.Flags().StringVarP(&confFormat, "format", "f", "", "配置文件格式 (json, yaml)")
}

func runConf(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// 确定格式
	format := confFormat
	if format == "" {
		if strings.HasSuffix(filename, ".json") {
			format = "json"
		} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
			format = "yaml"
		} else {
			return fmt.Errorf("无法确定配置文件格式，请使用 --format 参数指定")
		}
	}

	workDir, err := fileutil.GetWorkingDir()
	if err != nil {
		return err
	}

	// 创建日志目录
	logDir := filepath.Join(workDir, "log")
	if err := fileutil.EnsureDir(logDir); err != nil {
		return err
	}

	// 生成 cookie key
	cookieKey := web.MD5m16(time.Now().String() + filename)

	// 初始化模板引擎
	engine := template.NewEngine()

	// 注册模板
	if format == "json" {
		tmpl := `{
	"cookie": "{{ .CookieKey }}",
	"ip": "0.0.0.0",
	"port": 8080,
	"log": {
		"trace": "stdout",
		"info": "{{ .WorkDir }}/log/tigo-framework-info.log",
		"warning": "{{ .WorkDir }}/log/tigo-framework-warning.log",
		"error": "{{ .WorkDir }}/log/tigo-framework-error.log"
	}
}`
		if err := engine.Register("config", tmpl); err != nil {
			return err
		}
	} else {
		tmpl := `cookie: {{ .CookieKey }}
ip: 0.0.0.0
port: 8080
log:
  trace: stdout
  info: "{{ .WorkDir }}/log/tigo-framework-info.log"
  warning: "{{ .WorkDir }}/log/tigo-framework-warning.log"
  error: "{{ .WorkDir }}/log/tigo-framework-error.log"`
		if err := engine.Register("config", tmpl); err != nil {
			return err
		}
	}

	// 渲染配置文件
	data := template.TemplateData{
		"CookieKey": cookieKey,
		"WorkDir":   workDir,
	}
	configPath := filepath.Join(workDir, filename)
	if err := engine.RenderToFile("config", data, configPath); err != nil {
		return err
	}

	fmt.Printf("配置文件已生成: %s\n", configPath)

	// 更新 main.go 中的 ConfigPath
	if err := updateConfigPathInMain(workDir, filename); err != nil {
		fmt.Printf("警告: 无法自动更新 main.go 中的配置路径: %v\n", err)
		msg := fmt.Sprintf(`请手动添加: application.ConfigPath = "%s"`, filename)
		fmt.Println(msg)
		return nil
	}

	return nil
}

// updateConfigPathInMain 更新 main.go 中的配置路径
func updateConfigPathInMain(projectDir, filename string) error {
	mainPath := filepath.Join(projectDir, "main.go")

	content, err := fileutil.ReadFile(mainPath)
	if err != nil {
		return err
	}

	// 查找 application.Run() 并在其前面添加 ConfigPath
	configLine := `application.ConfigPath = "` + filename + `"` + "\n"

	if !strings.Contains(content, "ConfigPath") {
		content = strings.Replace(content,
			"application.Run()",
			configLine+"application.Run()",
			1)
	}

	return fileutil.WriteFile(mainPath, content)
}
