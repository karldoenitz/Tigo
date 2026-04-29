package cli

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// Prompter 交互式输入接口
type Prompter interface {
	AskString(question, defaultValue string) (string, error)
	AskSelect(question string, options []string) (string, error)
	Confirm(question string, defaultValue bool) (bool, error)
	AskMultiSelect(question string, options []string) ([]string, error)
	AskPassword(question string) (string, error)
}

// SurveyPrompter 使用 survey 库实现的 Prompter
type SurveyPrompter struct{}

// NewSurveyPrompter 创建新的 SurveyPrompter
func NewSurveyPrompter() *SurveyPrompter {
	return &SurveyPrompter{}
}

// AskString 询问字符串输入
func (p *SurveyPrompter) AskString(question, defaultValue string) (string, error) {
	var result string
	prompt := &survey.Input{
		Message: question,
		Default: defaultValue,
	}
	if err := survey.AskOne(prompt, &result); err != nil {
		return "", err
	}
	return result, nil
}

// AskSelect 询问选项选择
func (p *SurveyPrompter) AskSelect(question string, options []string) (string, error) {
	var result string
	prompt := &survey.Select{
		Message: question,
		Options: options,
	}
	if err := survey.AskOne(prompt, &result); err != nil {
		return "", err
	}
	return result, nil
}

// Confirm 确认操作
func (p *SurveyPrompter) Confirm(question string, defaultValue bool) (bool, error) {
	var result bool
	prompt := &survey.Confirm{
		Message: question,
		Default: defaultValue,
	}
	if err := survey.AskOne(prompt, &result); err != nil {
		return false, err
	}
	return result, nil
}

// AskMultiSelect 多选
func (p *SurveyPrompter) AskMultiSelect(question string, options []string) ([]string, error) {
	var result []string
	prompt := &survey.MultiSelect{
		Message: question,
		Options: options,
	}
	if err := survey.AskOne(prompt, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AskPassword 询问密码输入
func (p *SurveyPrompter) AskPassword(question string) (string, error) {
	var result string
	prompt := &survey.Password{
		Message: question,
	}
	if err := survey.AskOne(prompt, &result); err != nil {
		return "", err
	}
	return result, nil
}

// ValidateProjectName 验证项目名称
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("项目名称不能为空")
	}
	if len(name) > 64 {
		return fmt.Errorf("项目名称不能超过64个字符")
	}
	// 检查是否包含非法字符
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' || c == '_') {
			return fmt.Errorf("项目名称只能包含字母、数字、连字符和下划线")
		}
	}
	return nil
}

// SanitizeString 清理字符串，移除危险字符
func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "..", "")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "", "_")
	return s
}
