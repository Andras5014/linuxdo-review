package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Review   ReviewConfig   `yaml:"review"`
	SMTP     SMTPConfig     `yaml:"smtp"`
	OAuth    OAuthConfig    `yaml:"oauth"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

// ReviewConfig 审核配置
type ReviewConfig struct {
	MinVotes     int `yaml:"min_votes"`
	ApprovalRate int `yaml:"approval_rate"`
}

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

// OAuthConfig OAuth配置
type OAuthConfig struct {
	LinuxDo LinuxDoOAuthConfig `yaml:"linuxdo"`
}

// LinuxDoOAuthConfig Linux.do OAuth配置
type LinuxDoOAuthConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURI  string `yaml:"redirect_uri"`
}

var (
	cfg  *Config
	once sync.Once
)

// Load 加载配置文件
func Load(path string) (*Config, error) {
	var err error
	once.Do(func() {
		cfg = &Config{}
		var data []byte
		data, err = os.ReadFile(path)
		if err != nil {
			return
		}
		err = yaml.Unmarshal(data, cfg)
	})
	return cfg, err
}

// Get 获取全局配置
func Get() *Config {
	return cfg
}
