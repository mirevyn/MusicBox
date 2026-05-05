package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNormalizeBaseURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty", input: "", want: ""},
		{name: "trim whitespace", input: "  http://localhost:11434/  ", want: "http://localhost:11434"},
		{name: "trim trailing slashes", input: "https://example.com/v1///", want: "https://example.com/v1"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := normalizeBaseURL(tc.input); got != tc.want {
				t.Fatalf("normalizeBaseURL(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestNormalizeAIProvider(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		want  string
	}{
		{input: "ollama", want: "ollama"},
		{input: " OpenAI ", want: "openai"},
		{input: "auto", want: "auto"},
		{input: "unknown", want: ""},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			if got := normalizeAIProvider(tc.input); got != tc.want {
				t.Fatalf("normalizeAIProvider(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestValidateCustomAIBaseURLRejectsLocalAddresses(t *testing.T) {
	t.Parallel()

	tests := []string{
		"http://127.0.0.1:11434",
		"http://localhost:11434",
		"http://10.0.0.5:11434",
		"file:///tmp/socket",
	}

	for _, rawURL := range tests {
		rawURL := rawURL
		t.Run(rawURL, func(t *testing.T) {
			t.Parallel()
			err := ValidateCustomAIBaseURL(rawURL)
			if !errors.Is(err, ErrAIBaseURLNotAllowed) {
				t.Fatalf("ValidateCustomAIBaseURL(%q) error = %v, want ErrAIBaseURLNotAllowed", rawURL, err)
			}
		})
	}
}

func TestValidateCustomAIBaseURLAllowsPrivateWhenExplicitlyEnabled(t *testing.T) {
	t.Setenv("AI_ALLOW_PRIVATE_BASE_URLS", "true")

	if err := ValidateCustomAIBaseURL("http://127.0.0.1:11434"); err != nil {
		t.Fatalf("ValidateCustomAIBaseURL() error = %v, want nil", err)
	}
}

func TestLoadAITimeout(t *testing.T) {
	t.Parallel()

	key := "AI_REQUEST_TIMEOUT_SECONDS"
	original, existed := os.LookupEnv(key)
	defer restoreEnv(key, original, existed)

	if err := os.Setenv(key, "45"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	if got := loadAITimeout(); got != 45*time.Second {
		t.Fatalf("loadAITimeout() = %v, want %v", got, 45*time.Second)
	}

	if err := os.Setenv(key, "invalid"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	if got := loadAITimeout(); got != 300*time.Second {
		t.Fatalf("loadAITimeout() with invalid env = %v, want %v", got, 300*time.Second)
	}
}

func TestDetectProviderByExplicitV1Path(t *testing.T) {
	t.Parallel()

	srv := &AIService{
		baseURL:  "https://example.com/v1",
		provider: "auto",
		client:   http.DefaultClient,
	}

	got, err := srv.detectProvider()
	if err != nil {
		t.Fatalf("detectProvider() returned error: %v", err)
	}
	if got != "openai" {
		t.Fatalf("detectProvider() = %q, want %q", got, "openai")
	}
}

func TestDetectProviderByOllamaProbe(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/tags":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"models":[]}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	srv := &AIService{
		baseURL:  server.URL,
		provider: "auto",
		client:   server.Client(),
	}

	got, err := srv.detectProvider()
	if err != nil {
		t.Fatalf("detectProvider() returned error: %v", err)
	}
	if got != "ollama" {
		t.Fatalf("detectProvider() = %q, want %q", got, "ollama")
	}
}

func TestDetectProviderByOpenAIProbe(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/tags":
			http.NotFound(w, r)
		case "/v1/models":
			w.WriteHeader(http.StatusUnauthorized)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	srv := &AIService{
		baseURL:  server.URL,
		provider: "auto",
		client:   server.Client(),
	}

	got, err := srv.detectProvider()
	if err != nil {
		t.Fatalf("detectProvider() returned error: %v", err)
	}
	if got != "openai" {
		t.Fatalf("detectProvider() = %q, want %q", got, "openai")
	}
}

func TestDetectProviderFailsWhenUnknown(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	srv := &AIService{
		baseURL:  server.URL,
		provider: "auto",
		client:   server.Client(),
	}

	_, err := srv.detectProvider()
	if err == nil {
		t.Fatal("detectProvider() expected error, got nil")
	}
	if !errors.Is(err, ErrAIProviderDetection) {
		t.Fatalf("detectProvider() error = %v, want ErrAIProviderDetection", err)
	}
}

func restoreEnv(key, value string, existed bool) {
	if existed {
		_ = os.Setenv(key, value)
		return
	}
	_ = os.Unsetenv(key)
}
