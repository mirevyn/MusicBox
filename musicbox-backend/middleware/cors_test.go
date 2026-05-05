package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"musicbox-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func TestCORSSkipsWhenOriginsEmpty(t *testing.T) {
	t.Parallel()

	original := config.Conf
	t.Cleanup(func() {
		config.Conf = original
	})

	config.Conf = &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: nil,
		},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS())
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("unexpected CORS header %q when origins are empty", got)
	}
}
