package config

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestSplitAndTrim(t *testing.T) {
	t.Parallel()

	got := splitAndTrim(" http://a.com , , http://b.com ,, http://c.com ", ",")
	want := []string{"http://a.com", "http://b.com", "http://c.com"}

	if !slices.Equal(got, want) {
		t.Fatalf("splitAndTrim() = %#v, want %#v", got, want)
	}
}

func TestApplyEnvOverrides(t *testing.T) {
	t.Setenv("MUSICBOX_HTTP_ADDR", ":9000")
	t.Setenv("GIN_MODE", "release")
	t.Setenv("MUSICBOX_MYSQL_DSN", "root:secret@tcp(mysql:3306)/music_box")
	t.Setenv("MUSICBOX_JWT_SECRET", "jwt-secret")
	t.Setenv("MUSICBOX_JWT_EXPIRES_HOURS", "72")
	t.Setenv("MUSICBOX_UPLOAD_SONG_PATH", "custom/music")
	t.Setenv("MUSICBOX_UPLOAD_COVER_PATH", "custom/covers")
	t.Setenv("MUSICBOX_UPLOAD_LYRIC_PATH", "custom/lyrics")
	t.Setenv("MUSICBOX_UPLOAD_AVATAR_PATH", "custom/avatars")
	t.Setenv("MUSICBOX_CORS_ALLOWED_ORIGINS", "http://localhost:5173, https://app.example.com ")
	t.Setenv("AI_BASE_URL", "http://127.0.0.1:11434")
	t.Setenv("AI_MODEL", "qwen2.5:7b")
	t.Setenv("AI_PROVIDER", "ollama")
	t.Setenv("AI_API_KEY", "secret-key")
	t.Setenv("AI_REQUEST_TIMEOUT_SECONDS", "180")

	cfg := &Config{}
	applyEnvOverrides(cfg)

	if cfg.Server.HTTPAddr != ":9000" {
		t.Fatalf("Server.HTTPAddr = %q", cfg.Server.HTTPAddr)
	}
	if cfg.Server.GinMode != "release" {
		t.Fatalf("Server.GinMode = %q", cfg.Server.GinMode)
	}
	if cfg.MySQL.DSN != "root:secret@tcp(mysql:3306)/music_box" {
		t.Fatalf("MySQL.DSN = %q", cfg.MySQL.DSN)
	}
	if cfg.JWT.Secret != "jwt-secret" {
		t.Fatalf("JWT.Secret = %q", cfg.JWT.Secret)
	}
	if cfg.JWT.ExpiresTime != 72 {
		t.Fatalf("JWT.ExpiresTime = %d", cfg.JWT.ExpiresTime)
	}
	if cfg.Upload.SongPath != "custom/music" || cfg.Upload.CoverPath != "custom/covers" || cfg.Upload.LyricPath != "custom/lyrics" || cfg.Upload.AvatarPath != "custom/avatars" {
		t.Fatalf("upload paths not fully overridden: %#v", cfg.Upload)
	}

	wantOrigins := []string{"http://localhost:5173", "https://app.example.com"}
	if !slices.Equal(cfg.CORS.AllowedOrigins, wantOrigins) {
		t.Fatalf("CORS.AllowedOrigins = %#v, want %#v", cfg.CORS.AllowedOrigins, wantOrigins)
	}
	if cfg.AI.BaseURL != "http://127.0.0.1:11434" {
		t.Fatalf("AI.BaseURL = %q", cfg.AI.BaseURL)
	}
	if cfg.AI.Model != "qwen2.5:7b" {
		t.Fatalf("AI.Model = %q", cfg.AI.Model)
	}
	if cfg.AI.Provider != "ollama" {
		t.Fatalf("AI.Provider = %q", cfg.AI.Provider)
	}
	if cfg.AI.APIKey != "secret-key" {
		t.Fatalf("AI.APIKey = %q", cfg.AI.APIKey)
	}
	if cfg.AI.RequestTimeoutSeconds != 180 {
		t.Fatalf("AI.RequestTimeoutSeconds = %d", cfg.AI.RequestTimeoutSeconds)
	}
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	t.Run("missing mysql dsn panics", func(t *testing.T) {
		t.Parallel()
		expectPanic(t, func() {
			validateConfig(&Config{
				MySQL: MySQLConfig{DSN: ""},
				JWT:   JWTConfig{Secret: "secret"},
			})
		})
	})

	t.Run("missing jwt secret panics", func(t *testing.T) {
		t.Parallel()
		expectPanic(t, func() {
			validateConfig(&Config{
				MySQL: MySQLConfig{DSN: "dsn"},
				JWT:   JWTConfig{Secret: ""},
			})
		})
	})

		t.Run("complete config passes", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{
			MySQL: MySQLConfig{DSN: "dsn"},
			JWT:   JWTConfig{Secret: "secret"},
		}
		validateConfig(cfg)
		if cfg.Server.HTTPAddr != ":8000" {
			t.Fatalf("default Server.HTTPAddr = %q", cfg.Server.HTTPAddr)
		}
		if cfg.Server.GinMode != "debug" {
			t.Fatalf("default Server.GinMode = %q", cfg.Server.GinMode)
		}
		if cfg.JWT.ExpiresTime != 24 {
			t.Fatalf("default JWT.ExpiresTime = %d", cfg.JWT.ExpiresTime)
		}
		if cfg.AI.Model != "qwen2.5:7b" {
			t.Fatalf("default AI.Model = %q", cfg.AI.Model)
		}
		if cfg.AI.RequestTimeoutSeconds != 300 {
			t.Fatalf("default AI.RequestTimeoutSeconds = %d", cfg.AI.RequestTimeoutSeconds)
		}
	})
}

func TestConfigCandidatePathsIncludesExplicitPath(t *testing.T) {
	customPath := filepath.Join(t.TempDir(), "custom.yaml")
	t.Setenv("MUSICBOX_CONFIG_PATH", customPath)

	candidates := configCandidatePaths()
	if len(candidates) == 0 {
		t.Fatal("configCandidatePaths() returned no candidates")
	}
	if candidates[0] != customPath {
		t.Fatalf("first config candidate = %q, want %q", candidates[0], customPath)
	}
}

func TestConfigCandidatePathsPrefersLocalConfig(t *testing.T) {
	t.Setenv("MUSICBOX_CONFIG_PATH", "")

	candidates := configCandidatePaths()
	if len(candidates) < 2 {
		t.Fatalf("configCandidatePaths() = %#v, want at least local and example config", candidates)
	}
	if candidates[0] != "config/config.local.yaml" {
		t.Fatalf("first default config candidate = %q, want %q", candidates[0], "config/config.local.yaml")
	}
	if candidates[1] != "config/config.example.yaml" {
		t.Fatalf("second default config candidate = %q, want %q", candidates[1], "config/config.example.yaml")
	}
}

func TestReadConfigFileUsesExplicitPath(t *testing.T) {
	dir := t.TempDir()
	customPath := filepath.Join(dir, "musicbox-config.yaml")
	content := []byte("mysql:\n  dsn: explicit\njwt:\n  secret: explicit\n")
	if err := os.WriteFile(customPath, content, 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	t.Setenv("MUSICBOX_CONFIG_PATH", customPath)

	got, err := readConfigFile()
	if err != nil {
		t.Fatalf("readConfigFile() returned error: %v", err)
	}
	if string(got) != string(content) {
		t.Fatalf("readConfigFile() = %q, want %q", string(got), string(content))
	}
}

func expectPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic, got none")
		}
	}()
	fn()
}
