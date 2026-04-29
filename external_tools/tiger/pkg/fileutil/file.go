package fileutil

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnsureDir 确保目录存在
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	return nil
}

// Exists 检查文件或目录是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir 检查路径是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ReadFile 读取文件内容
func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	return string(content), nil
}

// WriteFile 写入文件
func WriteFile(path, content string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	return nil
}

// AppendFile 追加内容到文件
func AppendFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	return nil
}

// ReadLines 读取文件所有行
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	return lines, nil
}

// WriteLines 写入多行到文件
func WriteLines(path string, lines []string) error {
	content := strings.Join(lines, "\n")
	return WriteFile(path, content)
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("读取源文件失败: %w", err)
	}

	// 确保目标目录存在
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	if err := os.WriteFile(dst, content, 0644); err != nil {
		return fmt.Errorf("写入目标文件失败: %w", err)
	}
	return nil
}

// RemoveFile 删除文件
func RemoveFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// RemoveAll 删除目录或文件
func RemoveAll(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}
	return nil
}

// GetWorkingDir 获取当前工作目录
func GetWorkingDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取工作目录失败: %w", err)
	}
	return dir, nil
}
