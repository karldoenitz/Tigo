package cmd

import (
	"fmt"

	"github.com/karldoenitz/Tigo/web"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示 Tigo 版本信息",
	Long:  `显示当前安装的 Tigo 框架版本。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Tiger:", "2.0.0")
		fmt.Println("Tigo:", web.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
