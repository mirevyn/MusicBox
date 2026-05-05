package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"musicbox-backend/internal/config"
	"musicbox-backend/utils"

	"github.com/gin-gonic/gin"
)

func TestJWTAuth(t *testing.T) {
	setupJWTConfig(t)
	gin.SetMode(gin.TestMode)

	t.Run("rejects missing token", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/protected", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assertBusinessResponse(t, recorder, 401, "请求未携带Token，无权访问")
	})

	t.Run("rejects malformed token", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/protected", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "invalid-token")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assertBusinessResponse(t, recorder, 401, "Token格式错误")
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/protected", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assertBusinessResponse(t, recorder, 401, "Token无效或已过期")
	})

	t.Run("accepts valid token and writes claims to context", func(t *testing.T) {
		t.Parallel()

		token, _, err := utils.GenerateToken(7, "nebryx", "Admin")
		if err != nil {
			t.Fatalf("GenerateToken failed: %v", err)
		}

		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"userID":   c.GetUint("userID"),
				"username": c.GetString("username"),
				"role":     c.GetString("role"),
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
		}

		var payload map[string]any
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("unmarshal response failed: %v", err)
		}
		if got := uint(payload["userID"].(float64)); got != 7 {
			t.Fatalf("userID = %d, want %d", got, 7)
		}
		if got := payload["username"].(string); got != "nebryx" {
			t.Fatalf("username = %q, want %q", got, "nebryx")
		}
		if got := payload["role"].(string); got != "Admin" {
			t.Fatalf("role = %q, want %q", got, "Admin")
		}
	})
}

func TestAdminAuth(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("rejects missing role", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(AdminAuth())
		router.GET("/admin", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assertBusinessResponse(t, recorder, 403, "无法获取用户角色信息")
	})

	t.Run("rejects non admin role", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("role", "User")
			c.Next()
		})
		router.Use(AdminAuth())
		router.GET("/admin", func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assertBusinessResponse(t, recorder, 403, "权限不足，需要管理员身份")
	})

	t.Run("accepts admin role", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("role", "Admin")
			c.Next()
		})
		router.Use(AdminAuth())
		router.GET("/admin", func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
		}
		if recorder.Body.String() != "ok" {
			t.Fatalf("body = %q, want %q", recorder.Body.String(), "ok")
		}
	})
}

func setupJWTConfig(t *testing.T) {
	t.Helper()
	original := config.Conf
	config.Conf = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpiresTime: 24,
		},
	}
	t.Cleanup(func() {
		config.Conf = original
	})
}

func assertBusinessResponse(t *testing.T, recorder *httptest.ResponseRecorder, businessCode int, msg string) {
	t.Helper()

	if recorder.Code != businessCode {
		t.Fatalf("status code = %d, want %d", recorder.Code, businessCode)
	}

	var payload utils.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if payload.Code != businessCode {
		t.Fatalf("business code = %d, want %d", payload.Code, businessCode)
	}
	if payload.Msg != msg {
		t.Fatalf("msg = %q, want %q", payload.Msg, msg)
	}
}
