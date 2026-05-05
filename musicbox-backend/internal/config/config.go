package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
)

// 全局配置变量
var Conf *Config

// Config 结构体定义了应用的全部配置
type Config struct {
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	JWT    JWTConfig    `yaml:"jwt"`
	Upload UploadConfig `yaml:"upload"`
	CORS   CORSConfig   `yaml:"cors"`
	AI     AIConfig     `yaml:"ai"`
}

// UploadConfig 上传文件路径配置
type UploadConfig struct {
	SongPath   string `yaml:"song_path"`
	CoverPath  string `yaml:"cover_path"`
	LyricPath  string `yaml:"lyric_path"`
	AvatarPath string `yaml:"avatar_path"`
}

// MySQLConfig 数据库配置
type MySQLConfig struct {
	DSN string `yaml:"dsn"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpiresTime int    `yaml:"expires_time"`
}

// ServerConfig 服务运行配置
type ServerConfig struct {
	HTTPAddr string `yaml:"http_addr"`
	GinMode  string `yaml:"gin_mode"`
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

// AIConfig AI 服务配置
type AIConfig struct {
	BaseURL               string `yaml:"base_url"`
	Model                 string `yaml:"model"`
	Provider              string `yaml:"provider"`
	APIKey                string `yaml:"api_key"`
	RequestTimeoutSeconds int    `yaml:"request_timeout_seconds"`
}

// Init 函数用于加载配置文件
func Init() {
	data, err := readConfigFile()
	if err != nil {
		panic("无法读取配置文件: " + err.Error())
	}

	// 解析 YAML 到 Config 结构体
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic("无法解析配置文件: " + err.Error())
	}

	applyEnvOverrides(&cfg)
	validateConfig(&cfg)

	Conf = &cfg
}

func readConfigFile() ([]byte, error) {
	for _, path := range configCandidatePaths() {
		data, err := os.ReadFile(path)
		if err == nil {
			return data, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}
	return nil, os.ErrNotExist
}

func configCandidatePaths() []string {
	candidates := make([]string, 0, 5)
	seen := make(map[string]struct{})
	add := func(path string) {
		path = strings.TrimSpace(path)
		if path == "" {
			return
		}
		if _, ok := seen[path]; ok {
			return
		}
		seen[path] = struct{}{}
		candidates = append(candidates, path)
	}

	if explicit := os.Getenv("MUSICBOX_CONFIG_PATH"); explicit != "" {
		add(explicit)
	}

	add("config/config.local.yaml")
	add("config/config.example.yaml")

	if executable, err := os.Executable(); err == nil {
		execDir := filepath.Dir(executable)
		add(filepath.Join(execDir, "config", "config.local.yaml"))
		add(filepath.Join(execDir, "config", "config.example.yaml"))
	}

	return candidates
}

func applyEnvOverrides(cfg *Config) {
	if v := strings.TrimSpace(os.Getenv("MUSICBOX_HTTP_ADDR")); v != "" {
		cfg.Server.HTTPAddr = v
	}
	if v := strings.TrimSpace(os.Getenv("GIN_MODE")); v != "" {
		cfg.Server.GinMode = v
	}
	if dsn := os.Getenv("MUSICBOX_MYSQL_DSN"); dsn != "" {
		cfg.MySQL.DSN = dsn
	}
	if secret := os.Getenv("MUSICBOX_JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}
	if hoursStr := os.Getenv("MUSICBOX_JWT_EXPIRES_HOURS"); hoursStr != "" {
		if hours, err := strconv.Atoi(hoursStr); err == nil {
			cfg.JWT.ExpiresTime = hours
		}
	}
	if v := os.Getenv("MUSICBOX_UPLOAD_SONG_PATH"); v != "" {
		cfg.Upload.SongPath = v
	}
	if v := os.Getenv("MUSICBOX_UPLOAD_COVER_PATH"); v != "" {
		cfg.Upload.CoverPath = v
	}
	if v := os.Getenv("MUSICBOX_UPLOAD_LYRIC_PATH"); v != "" {
		cfg.Upload.LyricPath = v
	}
	if v := os.Getenv("MUSICBOX_UPLOAD_AVATAR_PATH"); v != "" {
		cfg.Upload.AvatarPath = v
	}
	if v, ok := os.LookupEnv("MUSICBOX_CORS_ALLOWED_ORIGINS"); ok {
		cfg.CORS.AllowedOrigins = splitAndTrim(v, ",")
	}
	if v := strings.TrimSpace(os.Getenv("AI_BASE_URL")); v != "" {
		cfg.AI.BaseURL = v
	}
	if v := strings.TrimSpace(os.Getenv("AI_MODEL")); v != "" {
		cfg.AI.Model = v
	}
	if v := strings.TrimSpace(os.Getenv("AI_PROVIDER")); v != "" {
		cfg.AI.Provider = v
	}
	if v := strings.TrimSpace(os.Getenv("AI_API_KEY")); v != "" {
		cfg.AI.APIKey = v
	}
	if v := strings.TrimSpace(os.Getenv("AI_REQUEST_TIMEOUT_SECONDS")); v != "" {
		if seconds, err := strconv.Atoi(v); err == nil {
			cfg.AI.RequestTimeoutSeconds = seconds
		}
	}
}

func validateConfig(cfg *Config) {
	if strings.TrimSpace(cfg.Server.HTTPAddr) == "" {
		cfg.Server.HTTPAddr = ":8000"
	}
	if strings.TrimSpace(cfg.Server.GinMode) == "" {
		cfg.Server.GinMode = "debug"
	}
	if strings.TrimSpace(cfg.MySQL.DSN) == "" {
		panic("MySQL DSN 未配置，请填写 config/config.local.yaml 或设置 MUSICBOX_MYSQL_DSN")
	}
	if strings.TrimSpace(cfg.JWT.Secret) == "" {
		panic("JWT Secret 未配置，请填写 config/config.local.yaml 或设置 MUSICBOX_JWT_SECRET")
	}
	if cfg.JWT.ExpiresTime <= 0 {
		cfg.JWT.ExpiresTime = 24
	}
	if strings.TrimSpace(cfg.AI.Model) == "" {
		cfg.AI.Model = "qwen2.5:7b"
	}
	if cfg.AI.RequestTimeoutSeconds <= 0 {
		cfg.AI.RequestTimeoutSeconds = 300
	}
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
