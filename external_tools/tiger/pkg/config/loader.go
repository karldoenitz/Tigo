package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config Tiger 配置
type Config struct {
	Version         string           `mapstructure:"version"`
	DefaultTemplate string           `mapstructure:"default_template"`
	Templates       []TemplateConfig `mapstructure:"templates"`
	TemplateRepos   []string         `mapstructure:"template_repositories"`
	ProjectDefaults ProjectDefaults  `mapstructure:"project_defaults"`
	Editor          EditorConfig     `mapstructure:"editor"`
	Git             GitConfig        `mapstructure:"git"`
}

// TemplateConfig 模板配置
type TemplateConfig struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

// ProjectDefaults 项目默认配置
type ProjectDefaults struct {
	Port          int    `mapstructure:"port"`
	IP            string `mapstructure:"ip"`
	EnableLogger  bool   `mapstructure:"enable_logger"`
	EnableSession bool   `mapstructure:"enable_session"`
}

// EditorConfig 编辑器配置
type EditorConfig struct {
	Format  string `mapstructure:"format"`
	Imports string `mapstructure:"imports"`
}

// GitConfig Git 配置
type GitConfig struct {
	AutoInit      bool   `mapstructure:"auto_init"`
	DefaultBranch string `mapstructure:"default_branch"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Version:         "2.0",
		DefaultTemplate: "basic",
		Templates:       []TemplateConfig{},
		TemplateRepos:   []string{},
		ProjectDefaults: ProjectDefaults{
			Port:          8080,
			IP:            "0.0.0.0",
			EnableLogger:  true,
			EnableSession: false,
		},
		Editor: EditorConfig{
			Format:  "gofmt",
			Imports: "goimports",
		},
		Git: GitConfig{
			AutoInit:      false,
			DefaultBranch: "main",
		},
	}
}

// Loader 配置加载器
type Loader struct {
	config *Config
}

// NewLoader 创建新的配置加载器
func NewLoader() *Loader {
	return &Loader{
		config: DefaultConfig(),
	}
}

// Load 从文件加载配置
func (l *Loader) Load(path string) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := v.Unmarshal(&l.config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// LoadFromEnv 从环境变量加载配置
func (l *Loader) LoadFromEnv() error {
	v := viper.New()
	v.SetEnvPrefix("TIGER")
	v.AutomaticEnv()

	// 映射环境变量到配置
	if port := v.GetInt("PORT"); port > 0 {
		l.config.ProjectDefaults.Port = port
	}
	if ip := v.GetString("IP"); ip != "" {
		l.config.ProjectDefaults.IP = ip
	}

	return nil
}

// Get 获取配置
func (l *Loader) Get() *Config {
	return l.config
}

// Save 保存配置到文件
func (l *Loader) Save(path string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	v := viper.New()
	v.Set("version", l.config.Version)
	v.Set("default_template", l.config.DefaultTemplate)
	v.Set("templates", l.config.Templates)
	v.Set("template_repositories", l.config.TemplateRepos)
	v.Set("project_defaults", l.config.ProjectDefaults)
	v.Set("editor", l.config.Editor)
	v.Set("git", l.config.Git)

	if err := v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// LoadDefault 加载默认配置
func LoadDefault() (*Config, error) {
	loader := NewLoader()

	// 尝试从用户主目录加载配置
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return loader.Get(), nil
	}

	configPath := filepath.Join(homeDir, ".tiger", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		if err := loader.Load(configPath); err != nil {
			return loader.Get(), err
		}
	}

	// 从环境变量加载
	_ = loader.LoadFromEnv()

	return loader.Get(), nil
}
