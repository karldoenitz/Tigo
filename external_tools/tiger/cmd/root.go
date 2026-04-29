package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tiger",
	Short: "Tigo 脚手架工具",
	Long: `Tiger 是 Tigo Web 框架的脚手架工具，
用于快速创建项目、添加 handler、生成配置文件等。`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
	},
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行失败: %v\n", err)
		os.Exit(1)
	}
}

// GetRootCmd 获取根命令
func GetRootCmd() *cobra.Command {
	return rootCmd
}
