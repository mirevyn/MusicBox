package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"musicbox-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func TestSwaggerDocUsesRequestHost(t *testing.T) {
	original := config.Conf
	t.Cleanup(func() {
		config.Conf = original
	})
	config.Conf = &config.Config{}

	gin.SetMode(gin.TestMode)
	router := InitRouter()

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	req.Host = "127.0.0.1:8000"
	req.Header.Set("X-Forwarded-Proto", "https")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
	}

	var payload struct {
		Host    string   `json:"host"`
		Schemes []string `json:"schemes"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal swagger doc failed: %v", err)
	}
	if payload.Host != "127.0.0.1:8000" {
		t.Fatalf("swagger host = %q, want %q", payload.Host, "127.0.0.1:8000")
	}
	if len(payload.Schemes) != 1 || payload.Schemes[0] != "https" {
		t.Fatalf("swagger schemes = %#v, want [https]", payload.Schemes)
	}
}
