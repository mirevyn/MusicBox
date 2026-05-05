package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"musicbox-backend/internal/service"
	"musicbox-backend/utils"

	"github.com/gin-gonic/gin"
)

func TestAIControllerGetStatusWithoutConfig(t *testing.T) {
	t.Setenv("AI_BASE_URL", "")
	t.Setenv("AI_MODEL", "")
	t.Setenv("AI_PROVIDER", "")
	t.Setenv("AI_API_KEY", "")

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/ai/status", nil)

	ctrl := &AIController{Service: service.NewAIService()}
	ctrl.GetStatus(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", recorder.Code, http.StatusOK)
	}

	var payload utils.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if payload.Code != 200 {
		t.Fatalf("business code = %d, want %d", payload.Code, 200)
	}
	if payload.Msg != "获取 AI 状态失败" {
		t.Fatalf("msg = %q, want %q", payload.Msg, "获取 AI 状态失败")
	}

	data, ok := payload.Data.(map[string]any)
	if !ok {
		t.Fatalf("payload data type = %T, want map[string]any", payload.Data)
	}
	if online, _ := data["online"].(bool); online {
		t.Fatal("online = true, want false")
	}
	if _, ok := data["error"].(string); !ok {
		t.Fatalf("error field missing in payload data: %#v", data)
	}
}

func TestAIControllerChatStreamRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/ai/chat/stream", bytes.NewBufferString("{invalid"))
	ctx.Request.Header.Set("Content-Type", "application/json")

	ctrl := &AIController{}
	ctrl.ChatStream(ctx)

	assertJSONResponse(t, recorder, 400, "参数错误")
}

func TestAIControllerChatStreamRejectsEmptyMessages(t *testing.T) {
	t.Parallel()

	body := AIChatRequest{
		Messages: []ChatMessageRequest{
			{Role: "user", Content: "   "},
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request failed: %v", err)
	}

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/ai/chat/stream", bytes.NewReader(raw))
	ctx.Request.Header.Set("Content-Type", "application/json")

	ctrl := &AIController{}
	ctrl.ChatStream(ctx)

	assertJSONResponse(t, recorder, 400, "消息不能为空")
}

func assertJSONResponse(t *testing.T, recorder *httptest.ResponseRecorder, businessCode int, msg string) {
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
