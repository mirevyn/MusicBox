package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"musicbox-backend/utils"

	"github.com/gin-gonic/gin"
)

func TestLoginRejectsInvalidPayload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
	}{
		{
			name: "invalid json",
			body: "{invalid",
		},
		{
			name: "missing password",
			body: `{"username":"admin"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gin.SetMode(gin.TestMode)
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(tt.body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			Login(ctx)

			assertBusinessMessage(t, recorder, 400, "参数错误:")
		})
	}
}

func TestFormatRegisterBindError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  string
		want string
	}{
		{
			name: "username required",
			err:  "Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag",
			want: "用户名不能为空",
		},
		{
			name: "username min",
			err:  "Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'min' tag",
			want: "用户名至少需要 3 个字符",
		},
		{
			name: "username max",
			err:  "Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'max' tag",
			want: "用户名长度不能超过 30 个字符",
		},
		{
			name: "password required",
			err:  "Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			want: "密码不能为空",
		},
		{
			name: "password min",
			err:  "Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag",
			want: "密码长度至少需要 6 位",
		},
		{
			name: "fallback",
			err:  "some other validation error",
			want: "注册参数错误",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := formatRegisterBindError(assertError(tt.err)); got != tt.want {
				t.Fatalf("formatRegisterBindError(%q) = %q, want %q", tt.err, got, tt.want)
			}
		})
	}
}

func assertBusinessMessage(t *testing.T, recorder *httptest.ResponseRecorder, code int, msgPrefix string) {
	t.Helper()

	if recorder.Code != code {
		t.Fatalf("status code = %d, want %d", recorder.Code, code)
	}

	var payload utils.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if payload.Code != code {
		t.Fatalf("business code = %d, want %d", payload.Code, code)
	}
	if len(payload.Msg) < len(msgPrefix) || payload.Msg[:len(msgPrefix)] != msgPrefix {
		t.Fatalf("msg = %q, want prefix %q", payload.Msg, msgPrefix)
	}
}

type assertError string

func (e assertError) Error() string {
	return string(e)
}
