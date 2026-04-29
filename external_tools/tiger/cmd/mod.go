package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/karldoenitz/Tigo/external_tools/tiger/pkg/fileutil"
	"github.com/spf13/cobra"
)

var modCmd = &cobra.Command{
	Use:   "mod",
	Short: "执行 go mod 命令",
	Long:  `执行 go mod init, go mod tidy 和 go mod vendor 命令。`,
	RunE:  runMod,
}

func init() {
	rootCmd.AddCommand(modCmd)
}

func runMod(cmd *cobra.Command, args []string) error {
	workDir, err := fileutil.GetWorkingDir()
	if err != nil {
		return err
	}

	// 获取项目名
	splitPath := strings.Split(workDir, string(filepath.Separator))
	proName := splitPath[len(splitPath)-1]

	// 设置环境变量
	_ = exec.Command("/bin/sh", "-c", "export GO111MODULE=on").Run()

	// 执行 go mod init
	fmt.Println("执行: go mod init", proName)
	cmdInit := exec.Command("go", "mod", "init", proName)
	cmdInit.Stdout = os.Stdout
	cmdInit.Stderr = os.Stderr
	if err := cmdInit.Run(); err != nil {
		fmt.Printf("go mod init 失败: %v\n", err)
	}

	// 执行 go mod tidy
	fmt.Println("执行: go mod tidy")
	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Stdout = os.Stdout
	cmdTidy.Stderr = os.Stderr
	if err := cmdTidy.Run(); err != nil {
		fmt.Printf("go mod tidy 失败: %v\n", err)
	}

	// 执行 go mod vendor
	fmt.Println("执行: go mod vendor")
	cmdVendor := exec.Command("go", "mod", "vendor")
	cmdVendor.Stdout = os.Stdout
	cmdVendor.Stderr = os.Stderr
	if err := cmdVendor.Run(); err != nil {
		fmt.Printf("go mod vendor 失败: %v\n", err)
	}

	fmt.Println("Go mod 完成!")
	return nil
}
