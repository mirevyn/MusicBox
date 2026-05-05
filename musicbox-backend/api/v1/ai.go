package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"musicbox-backend/internal/service"
	"musicbox-backend/utils"

	"github.com/gin-gonic/gin"
)

type AIController struct {
	Service *service.AIService
}

type ChatMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIChatRequest struct {
	Messages []ChatMessageRequest `json:"messages" binding:"required"`
	BaseURL  string               `json:"baseUrl"`
	Model    string               `json:"model"`
	ApiKey   string               `json:"apiKey"`
	Provider string               `json:"provider"`
}

type AIConfigRequest struct {
	BaseURL  string `json:"baseUrl"`
	Model    string `json:"model"`
	ApiKey   string `json:"apiKey"`
	Provider string `json:"provider"`
}

const (
	maxAIChatMessages      = 20
	maxAIMessageRunes      = 2000
	maxAIChatTotalRunes    = 12000
	defaultAIErrorResponse = "获取 AI 状态失败"
)

func NewAIController() *AIController {
	return &AIController{
		Service: service.NewAIService(),
	}
}

// @Summary 获取 AI 状态
// @Description 检查当前 AI 服务及模型的可用性
// @Tags AI 助手
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body AIConfigRequest false "AI 服务配置"
// @Success 200 {object} utils.Response "成功"
// @Router /ai/status [post]
func (ctrl *AIController) GetStatus(c *gin.Context) {
	req, ok := bindAIConfigRequest(c)
	if !ok {
		return
	}

	srv := ctrl.Service
	if req.BaseURL != "" || req.Model != "" || req.ApiKey != "" || req.Provider != "" {
		if err := service.ValidateCustomAIBaseURL(req.BaseURL); err != nil {
			utils.Result(c, 400, nil, err.Error())
			return
		}
		srv = srv.WithConfig(req.BaseURL, req.Model, req.ApiKey, req.Provider)
	}

	online, activeModel, models, err := srv.Status()
	if err != nil {
		utils.Result(c, 200, gin.H{
			"online":      false,
			"model":       activeModel,
			"models":      []string{},
			"providerKey": srv.ProviderKey(),
			"provider":    srv.ProviderName(),
			"error":       err.Error(),
			"baseUrl":     srv.BaseURL(), // 方便检查实际使用的地址
		}, defaultAIErrorResponse)
		return
	}

	utils.Result(c, 200, gin.H{
		"online":      online,
		"model":       activeModel,
		"models":      models,
		"providerKey": srv.ProviderKey(),
		"provider":    srv.ProviderName(),
		"baseUrl":     srv.BaseURL(),
	}, "成功")
}

// @Summary AI 聊天 (流式)
// @Description 与 AI 音乐助手进行流式对话，支持实时文本输出和最后返回推荐歌曲
// @Tags AI 助手
// @Security ApiKeyAuth
// @Accept json
// @Produce text/event-stream
// @Param request body AIChatRequest true "聊天请求参数"
// @Success 200 {string} string "SSE 事件流"
// @Router /ai/chat/stream [post]
func (ctrl *AIController) ChatStream(c *gin.Context) {
	var req AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(c, 400, nil, "参数错误")
		return
	}
	if len(req.Messages) > maxAIChatMessages {
		utils.Result(c, 400, nil, "消息数量过多")
		return
	}
	if err := service.ValidateCustomAIBaseURL(req.BaseURL); err != nil {
		utils.Result(c, 400, nil, err.Error())
		return
	}

	messages := make([]service.AIChatMessage, 0, len(req.Messages))
	totalRunes := 0
	for _, item := range req.Messages {
		role, ok := normalizeAIChatRole(item.Role)
		if !ok {
			utils.Result(c, 400, nil, "消息角色不合法")
			return
		}
		content := strings.TrimSpace(item.Content)
		if content == "" {
			continue
		}
		contentRunes := len([]rune(content))
		if contentRunes > maxAIMessageRunes {
			utils.Result(c, 400, nil, "单条消息过长")
			return
		}
		totalRunes += contentRunes
		if totalRunes > maxAIChatTotalRunes {
			utils.Result(c, 400, nil, "消息内容过长")
			return
		}
		messages = append(messages, service.AIChatMessage{
			Role:    role,
			Content: content,
		})
	}

	if len(messages) == 0 {
		utils.Result(c, 400, nil, "消息不能为空")
		return
	}

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		utils.Result(c, 500, nil, "当前服务不支持流式输出")
		return
	}

	uid := getOptionalAuthUserID(c)

	srv := ctrl.Service
	if req.BaseURL != "" || req.Model != "" || req.ApiKey != "" || req.Provider != "" {
		srv = srv.WithConfig(req.BaseURL, req.Model, req.ApiKey, req.Provider)
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache, no-transform")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	var recommendedSongs []map[string]interface{}
	flushEvent := func(event string, payload interface{}) error {
		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if _, err := c.Writer.Write([]byte("event: " + event + "\n")); err != nil {
			return err
		}
		if _, err := c.Writer.Write([]byte("data: " + string(body) + "\n\n")); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	candidateSongs, _, err := srv.StreamChat(messages, uid, func(chunk string) error {
		return flushEvent("chunk", gin.H{"content": chunk})
	})
	if err != nil {
		if errors.Is(err, service.ErrAIUnavailable) {
			_ = flushEvent("error", gin.H{"message": err.Error()})
			return
		}
		_ = flushEvent("error", gin.H{"message": err.Error()})
		return
	}

	for _, song := range candidateSongs {
		recommendedSongs = append(recommendedSongs, gin.H{
			"id":       song.ID,
			"songId":   song.ID,
			"title":    song.Title,
			"artist":   song.Artist,
			"album":    song.Album,
			"duration": song.Duration,
			"coverUrl": song.CoverUrl,
			"fileUrl":  song.FileUrl,
			"lyricUrl": song.LyricUrl,
		})
	}
	_ = flushEvent("meta", gin.H{
		"model":       srv.ModelName(),
		"provider":    srv.ProviderName(),
		"providerKey": srv.ProviderKey(),
		"songs":       recommendedSongs,
	})
	_ = flushEvent("done", gin.H{"ok": true})
}

func bindAIConfigRequest(c *gin.Context) (AIConfigRequest, bool) {
	if c.Request.Method == http.MethodGet {
		if strings.TrimSpace(c.Query("apiKey")) != "" {
			utils.Result(c, 400, nil, "API Key 不允许通过 URL 参数传递")
			return AIConfigRequest{}, false
		}
		return AIConfigRequest{
			BaseURL:  strings.TrimSpace(c.Query("baseUrl")),
			Model:    strings.TrimSpace(c.Query("model")),
			Provider: strings.TrimSpace(c.Query("provider")),
		}, true
	}

	var req AIConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(c, 400, nil, "参数错误")
		return AIConfigRequest{}, false
	}
	req.BaseURL = strings.TrimSpace(req.BaseURL)
	req.Model = strings.TrimSpace(req.Model)
	req.ApiKey = strings.TrimSpace(req.ApiKey)
	req.Provider = strings.TrimSpace(req.Provider)
	return req, true
}

func normalizeAIChatRole(role string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "user":
		return "user", true
	case "assistant":
		return "assistant", true
	default:
		return "", false
	}
}
