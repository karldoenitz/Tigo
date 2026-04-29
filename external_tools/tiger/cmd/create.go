package cmd

import (
	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/cli"
	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/project"
	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/template"
	"github.com/spf13/cobra"
)

var (
	createTemplate    string
	createPort        int
	createIP          string
	createLogger      bool
	createSession     bool
	createTests       bool
	createGit         bool
	createInteractive bool
)

var createCmd = &cobra.Command{
	Use:   "create <project_name>",
	Short: "创建一个新的 Tigo 项目",
	Long: `创建一个新的 Tigo 项目。
支持多种模板类型，可以使用 --interactive 参数进行交互式创建。`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createTemplate, "template", "t", "", "项目模板 (basic, restful, websocket, graphql, full)")
	createCmd.Flags().IntVarP(&createPort, "port", "p", 0, "服务端口")
	createCmd.Flags().StringVarP(&createIP, "ip", "i", "", "服务 IP 地址")
	createCmd.Flags().BoolVar(&createLogger, "enable-logger", false, "启用日志")
	createCmd.Flags().BoolVar(&createSession, "enable-session", false, "启用 Session")
	createCmd.Flags().BoolVar(&createTests, "include-tests", false, "包含测试文件")
	createCmd.Flags().BoolVar(&createGit, "git-init", false, "初始化 Git")
	createCmd.Flags().BoolVarP(&createInteractive, "interactive", "I", false, "交互式创建")
}

func runCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	// 初始化模板引擎
	engine := template.NewEngine()

	// 创建项目创建器
	prompter := cli.NewSurveyPrompter()
	creator := project.NewCreator(engine, prompter)

	// 创建选项
	opts := project.CreateOptions{
		Name:          name,
		Template:      createTemplate,
		Port:          createPort,
		IP:            createIP,
		EnableLogger:  createLogger,
		EnableSession: createSession,
		IncludeTests:  createTests,
		GitInit:       createGit,
		Interactive:   createInteractive,
	}

	// 执行创建
	return creator.Create(opts)
}
