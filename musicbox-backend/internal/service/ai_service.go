package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"musicbox-backend/internal/config"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"musicbox-backend/internal/global"
	"musicbox-backend/internal/model"

	"github.com/gin-gonic/gin"
)

var ErrAIUnavailable = errors.New("ai service unavailable")
var ErrAIConfigMissing = errors.New("ai service not configured")
var ErrAIProviderDetection = errors.New("ai provider detection failed")
var ErrAIBaseURLNotAllowed = errors.New("ai base url not allowed")

type AIService struct {
	baseURL  string
	model    string
	apiKey   string
	provider string
	detected string
	client   *http.Client
}

type AIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []AIChatMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaChatResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

type ollamaStreamResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done  bool   `json:"done"`
	Error string `json:"error"`
}

type ollamaTagsResponse struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}

type AIChatContext struct {
	Messages     []AIChatMessage
	Candidates   []model.Song
	DirectReply  string
	RequestTopic string
}

type AIIntentAnalysis struct {
	IsRecommendation bool     `json:"isRecommendation"`
	IsFollowup       bool     `json:"isFollowup"`
	Keywords         []string `json:"keywords"`
	Mood             string   `json:"mood"`
	Language         string   `json:"language"`
}

type scoredSongCandidate struct {
	Song  model.Song
	Score int
}

type userPreferenceProfile struct {
	SongWeights   map[uint]int
	ArtistWeights map[string]int
	AlbumWeights  map[string]int
}

type aiSongSelectionResult struct {
	SongIDs []uint `json:"songIds"`
}

func NewAIService() *AIService {
	baseURL := loadAIBaseURL()
	modelName := loadAIModel()

	return &AIService{
		baseURL:  normalizeBaseURL(baseURL),
		model:    modelName,
		apiKey:   loadAIAPIKey(),
		provider: loadAIProvider(),
		detected: "",
		client: &http.Client{
			Timeout: loadAITimeout(),
		},
	}
}

// normalizeBaseURL 统一清洗外部传入的基础地址，避免尾部斜杠影响拼接。
func normalizeBaseURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return ""
	}
	return strings.TrimRight(url, "/")
}

// ValidateCustomAIBaseURL 校验由前端传入的 AI 服务地址，避免把后端变成任意内网探测代理。
func ValidateCustomAIBaseURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return fmt.Errorf("%w: AI 服务地址格式不正确", ErrAIBaseURLNotAllowed)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("%w: AI 服务地址只允许 http 或 https", ErrAIBaseURLNotAllowed)
	}
	if parsed.User != nil {
		return fmt.Errorf("%w: AI 服务地址不能包含用户名或密码", ErrAIBaseURLNotAllowed)
	}

	if allowPrivateCustomAIBaseURLs() {
		return nil
	}

	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if host == "" || host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return fmt.Errorf("%w: 不允许连接本机或内网地址", ErrAIBaseURLNotAllowed)
	}

	if ip := net.ParseIP(host); ip != nil {
		if isBlockedAIHostIP(ip) {
			return fmt.Errorf("%w: 不允许连接本机或内网地址", ErrAIBaseURLNotAllowed)
		}
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil || len(addrs) == 0 {
		return fmt.Errorf("%w: 无法解析 AI 服务地址", ErrAIBaseURLNotAllowed)
	}
	for _, addr := range addrs {
		if isBlockedAIHostIP(addr.IP) {
			return fmt.Errorf("%w: 不允许连接本机或内网地址", ErrAIBaseURLNotAllowed)
		}
	}
	return nil
}

func allowPrivateCustomAIBaseURLs() bool {
	value := strings.TrimSpace(os.Getenv("MUSICBOX_AI_ALLOW_PRIVATE_BASE_URLS"))
	if value == "" {
		value = strings.TrimSpace(os.Getenv("AI_ALLOW_PRIVATE_BASE_URLS"))
	}
	return strings.EqualFold(value, "true") || value == "1"
}

func isBlockedAIHostIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return true
	}
	addr = addr.Unmap()
	if !addr.IsGlobalUnicast() || addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast() || addr.IsMulticast() || addr.IsUnspecified() {
		return true
	}

	for _, cidr := range []string{
		"0.0.0.0/8",
		"100.64.0.0/10",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"224.0.0.0/4",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	} {
		if netip.MustParsePrefix(cidr).Contains(addr) {
			return true
		}
	}
	return false
}

func loadAITimeout() time.Duration {
	timeoutSeconds := 300
	if cfg := config.Conf; cfg != nil && cfg.AI.RequestTimeoutSeconds > 0 {
		timeoutSeconds = cfg.AI.RequestTimeoutSeconds
	}
	if raw := strings.TrimSpace(os.Getenv("AI_REQUEST_TIMEOUT_SECONDS")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			timeoutSeconds = parsed
		}
	}
	return time.Duration(timeoutSeconds) * time.Second
}

func loadAIBaseURL() string {
	if cfg := config.Conf; cfg != nil && strings.TrimSpace(cfg.AI.BaseURL) != "" {
		return cfg.AI.BaseURL
	}
	if value := strings.TrimSpace(os.Getenv("AI_BASE_URL")); value != "" {
		return value
	}
	return ""
}

func loadAIModel() string {
	if cfg := config.Conf; cfg != nil && strings.TrimSpace(cfg.AI.Model) != "" {
		return cfg.AI.Model
	}
	if value := strings.TrimSpace(os.Getenv("AI_MODEL")); value != "" {
		return value
	}
	return "qwen2.5:7b"
}

func loadAIProvider() string {
	value := ""
	if cfg := config.Conf; cfg != nil {
		value = cfg.AI.Provider
	}
	if envValue := strings.TrimSpace(os.Getenv("AI_PROVIDER")); envValue != "" {
		value = envValue
	}
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "ollama", "openai", "auto":
		return value
	default:
		return "auto"
	}
}

func loadAIAPIKey() string {
	if cfg := config.Conf; cfg != nil && strings.TrimSpace(cfg.AI.APIKey) != "" {
		return cfg.AI.APIKey
	}
	return strings.TrimSpace(os.Getenv("AI_API_KEY"))
}

// WithConfig 会创建一个新的 AIService 实例，并使用新的 URL、Model 和 apiKey 配置进行替换。
func (s *AIService) WithConfig(baseURL, model, apiKey, provider string) *AIService {
	ns := &AIService{
		baseURL:  s.baseURL,
		model:    s.model,
		apiKey:   s.apiKey,
		provider: s.provider,
		detected: s.detected,
		client:   s.client,
	}
	if strings.TrimSpace(baseURL) != "" {
		ns.baseURL = normalizeBaseURL(baseURL)
		ns.detected = ""
	}
	if strings.TrimSpace(model) != "" {
		ns.model = strings.TrimSpace(model)
	}
	if strings.TrimSpace(apiKey) != "" {
		ns.apiKey = strings.TrimSpace(apiKey)
		ns.detected = ""
	}
	if normalizedProvider := normalizeAIProvider(provider); normalizedProvider != "" {
		ns.provider = normalizedProvider
		ns.detected = ""
	}
	return ns
}

func normalizeAIProvider(provider string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "ollama":
		return "ollama"
	case "openai":
		return "openai"
	case "auto":
		return "auto"
	default:
		return ""
	}
}

func (s *AIService) requestChatCompletion(messages []AIChatMessage, stream bool) (string, error) {
	if err := s.ensureConfigured(); err != nil {
		return "", err
	}
	provider, err := s.effectiveProvider()
	if err != nil {
		return "", err
	}
	if provider == "ollama" {
		return s.requestOllamaCompletion(messages, stream)
	}
	return s.requestOpenAICompletion(messages, stream)
}

func (s *AIService) effectiveProvider() (string, error) {
	switch s.provider {
	case "ollama":
		return "ollama", nil
	case "openai":
		return "openai", nil
	}
	if s.detected != "" {
		return s.detected, nil
	}

	detected, err := s.detectProvider()
	if err != nil {
		return "", err
	}
	s.detected = detected
	return detected, nil
}

func (s *AIService) detectProvider() (string, error) {
	baseURL := normalizeBaseURL(s.baseURL)
	if baseURL == "" {
		return "", ErrAIConfigMissing
	}

	normalized := strings.ToLower(baseURL)
	if strings.HasSuffix(normalized, "/v1") || strings.Contains(normalized, "/v1/") {
		return "openai", nil
	}

	if ok := s.probeOllama(baseURL); ok {
		return "ollama", nil
	}
	if ok := s.probeOpenAI(); ok {
		return "openai", nil
	}

	return "", fmt.Errorf("%w: 无法自动识别当前 AI 服务类型，请在配置中心明确选择 Ollama 或 OpenAI 兼容", ErrAIProviderDetection)
}

func (s *AIService) probeOllama(baseURL string) bool {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 400
}

func (s *AIService) probeOpenAI() bool {
	req, err := http.NewRequest(http.MethodGet, s.getOpenAIUrl()+"/models", nil)
	if err != nil {
		return false
	}
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 400 || resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden
}

func (s *AIService) requestOllamaCompletion(messages []AIChatMessage, stream bool) (string, error) {
	payload, err := json.Marshal(ollamaChatRequest{
		Model:    s.model,
		Messages: messages,
		Stream:   stream,
	})
	if err != nil {
		return "", err
	}
	var errs []string
	for _, baseURL := range s.ollamaCandidateBaseURLs() {
		req, err := http.NewRequest(http.MethodPost, baseURL+"/api/chat", bytes.NewReader(payload))
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s -> %v", baseURL, err))
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.client.Do(req)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s -> %v", baseURL, err))
			continue
		}

		if resp.StatusCode >= 400 {
			_ = resp.Body.Close()
			errs = append(errs, fmt.Sprintf("%s -> server returned %d", baseURL, resp.StatusCode))
			continue
		}

		var result ollamaChatResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		_ = resp.Body.Close()
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s -> %v", baseURL, err))
			continue
		}

		content := strings.TrimSpace(result.Message.Content)
		if content == "" {
			errs = append(errs, fmt.Sprintf("%s -> empty content", baseURL))
			continue
		}

		s.baseURL = baseURL
		return content, nil
	}

	return "", fmt.Errorf("%w: %s", ErrAIUnavailable, strings.Join(errs, "; "))
}

func (s *AIService) requestOpenAICompletion(messages []AIChatMessage, stream bool) (string, error) {
	url := s.getOpenAIUrl()
	payload, err := json.Marshal(gin.H{
		"model":    s.model,
		"messages": messages,
		"stream":   stream,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrAIUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("%w: OpenAI API returned %d", ErrAIUnavailable, resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", ErrAIUnavailable
	}

	return result.Choices[0].Message.Content, nil
}

// StreamChat 代理 Ollama 或 OpenAI 的流式输出
func (s *AIService) StreamChat(history []AIChatMessage, userID uint, onChunk func(string) error) ([]model.Song, string, error) {
	if err := s.ensureConfigured(); err != nil {
		return nil, s.model, err
	}
	context, err := s.prepareChatContext(history, userID)
	if err != nil {
		return nil, s.model, err
	}
	if context.DirectReply != "" {
		if err := onChunk(context.DirectReply); err != nil {
			return nil, s.model, err
		}
		return nil, s.model, nil
	}

	provider, err := s.effectiveProvider()
	if err != nil {
		return nil, s.model, err
	}
	if provider == "ollama" {
		return s.streamOllamaChat(context.Messages, onChunk, context.Candidates)
	}
	return s.streamOpenAIChat(context.Messages, onChunk, context.Candidates)
}

func (s *AIService) streamOllamaChat(messages []AIChatMessage, onChunk func(string) error, candidates []model.Song) ([]model.Song, string, error) {
	payload, err := json.Marshal(ollamaChatRequest{
		Model:    s.model,
		Messages: messages,
		Stream:   true,
	})
	if err != nil {
		return nil, s.model, err
	}
	var errs []string
	for _, baseURL := range s.ollamaCandidateBaseURLs() {
		req, err := http.NewRequest(http.MethodPost, baseURL+"/api/chat", bytes.NewReader(payload))
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s -> %v", baseURL, err))
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.client.Do(req)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s -> %v", baseURL, err))
			continue
		}

		if resp.StatusCode >= 400 {
			_ = resp.Body.Close()
			errs = append(errs, fmt.Sprintf("%s -> server returned %d", baseURL, resp.StatusCode))
			continue
		}

		scanner := bufio.NewScanner(resp.Body)
		streamErr := error(nil)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			var chunk ollamaStreamResponse
			if err := json.Unmarshal([]byte(line), &chunk); err != nil {
				streamErr = err
				break
			}
			if chunk.Error != "" {
				streamErr = fmt.Errorf("%w: %s", ErrAIUnavailable, chunk.Error)
				break
			}

			piece := chunk.Message.Content
			if piece != "" {
				if err := onChunk(piece); err != nil {
					_ = resp.Body.Close()
					return nil, s.model, err
				}
			}

			if chunk.Done {
				break
			}
		}

		if streamErr == nil {
			streamErr = scanner.Err()
		}
		_ = resp.Body.Close()
		if streamErr != nil {
			errs = append(errs, fmt.Sprintf("%s -> %v", baseURL, streamErr))
			continue
		}

		s.baseURL = baseURL
		return candidates, s.model, nil
	}

	return nil, s.model, fmt.Errorf("%w: %s", ErrAIUnavailable, strings.Join(errs, "; "))
}

func (s *AIService) streamOpenAIChat(messages []AIChatMessage, onChunk func(string) error, candidates []model.Song) ([]model.Song, string, error) {
	url := s.baseURL
	if !strings.HasSuffix(url, "/v1") && !strings.Contains(url, "/v1/") {
		url = strings.TrimRight(url, "/") + "/v1"
	}

	payload, err := json.Marshal(gin.H{
		"model":    s.model,
		"messages": messages,
		"stream":   true,
	})
	if err != nil {
		return nil, s.model, err
	}

	req, err := http.NewRequest(http.MethodPost, url+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, s.model, err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, s.model, fmt.Errorf("%w: %v", ErrAIUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, s.model, fmt.Errorf("%w: OpenAI API returned %d", ErrAIUnavailable, resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, s.model, err
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				if err := onChunk(content); err != nil {
					return nil, s.model, err
				}
			}
		}
	}

	return candidates, s.model, nil
}

// Status 检查 AI 服务是否在线，并确认目标模型是否已存在。
func (s *AIService) Status() (bool, string, []string, error) {
	if err := s.ensureConfigured(); err != nil {
		return false, s.model, []string{}, err
	}
	provider, err := s.effectiveProvider()
	if err != nil {
		return false, s.model, []string{}, err
	}
	if provider == "ollama" {
		return s.statusOllama()
	}
	return s.statusOpenAI()
}

func (s *AIService) ensureConfigured() error {
	if strings.TrimSpace(s.baseURL) != "" {
		return nil
	}
	return fmt.Errorf("%w: 请先配置 AI_BASE_URL，或在前端 AI 面板里保存一个连接方案", ErrAIConfigMissing)
}

func (s *AIService) statusOllama() (bool, string, []string, error) {
	var errs []string
	for _, baseURL := range s.ollamaCandidateBaseURLs() {
		req, err := http.NewRequest(http.MethodGet, baseURL+"/api/tags", nil)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s -> %v", baseURL, err))
			continue
		}

		resp, err := s.client.Do(req)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s -> %v", baseURL, err))
			continue
		}

		if resp.StatusCode >= 400 {
			errs = append(errs, fmt.Sprintf("%s -> %s", baseURL, resp.Status))
			_ = resp.Body.Close()
			continue
		}

		var tags ollamaTagsResponse
		err = json.NewDecoder(resp.Body).Decode(&tags)
		_ = resp.Body.Close()
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s -> %v", baseURL, err))
			continue
		}

		models := make([]string, 0, len(tags.Models))
		modelOk := false
		for _, item := range tags.Models {
			models = append(models, item.Name)
			if strings.EqualFold(strings.TrimSpace(item.Name), s.model) {
				modelOk = true
			}
		}

		s.baseURL = baseURL
		return modelOk, s.model, models, nil
	}

	return false, s.model, nil, fmt.Errorf("连接 Ollama 失败: %s", strings.Join(errs, "; "))
}

func (s *AIService) statusOpenAI() (bool, string, []string, error) {
	url := s.getOpenAIUrl()
	req, err := http.NewRequest(http.MethodGet, url+"/models", nil)
	if err != nil {
		return false, s.model, nil, err
	}
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return false, s.model, nil, fmt.Errorf("连接 OpenAI API 失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return false, s.model, nil, fmt.Errorf("API 密钥无效或无权限 (HTTP %d)", resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		return false, s.model, nil, fmt.Errorf("OpenAI API 响应错误: %s", resp.Status)
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return true, s.model, []string{s.model}, nil
	}

	models := make([]string, 0, len(result.Data))
	modelOk := strings.TrimSpace(s.model) == ""
	for _, item := range result.Data {
		name := strings.TrimSpace(item.ID)
		if name == "" {
			continue
		}
		models = append(models, name)
		if strings.EqualFold(name, s.model) {
			modelOk = true
		}
	}
	if len(models) == 0 {
		models = []string{s.model}
		modelOk = true
	}

	return modelOk, s.model, models, nil
}

func (s *AIService) getOpenAIUrl() string {
	url := s.baseURL
	if !strings.HasSuffix(url, "/v1") && !strings.Contains(url, "/v1/") {
		url = strings.TrimRight(url, "/") + "/v1"
	}
	return url
}

func (s *AIService) BaseURL() string {
	return s.baseURL
}

func (s *AIService) ModelName() string {
	return s.model
}

func (s *AIService) ProviderKey() string {
	provider, err := s.effectiveProvider()
	if err != nil {
		if s.provider != "auto" {
			return s.provider
		}
		return "auto"
	}
	if provider == "ollama" {
		return "ollama"
	}
	return "openai"
}

func (s *AIService) ProviderName() string {
	if s.ProviderKey() == "ollama" {
		return "Ollama"
	}
	if s.ProviderKey() == "openai" {
		return "OpenAI-compatible"
	}
	return "自动识别"
}

func (s *AIService) isOllamaNative() bool {
	provider, err := s.effectiveProvider()
	return err == nil && provider == "ollama"
}

func (s *AIService) ollamaCandidateBaseURLs() []string {
	baseURL := normalizeBaseURL(s.baseURL)
	if baseURL == "" {
		return nil
	}
	return []string{baseURL}
}

func (s *AIService) prepareChatContext(history []AIChatMessage, userID uint) (*AIChatContext, error) {
	prompt, intent := s.resolveRecommendationIntent(history)
	var (
		candidates []model.Song
		err        error
	)

	if intent.IsRecommendation {
		candidates, err = s.getRecommendationCandidates(prompt, intent, userID, 12)
	} else {
		candidates = nil
	}
	if err != nil {
		return nil, err
	}
	if intent.IsRecommendation && len(candidates) == 0 {
		return &AIChatContext{
			DirectReply:  fmt.Sprintf("当前曲库缺少与“%s”直接匹配的歌曲。你可以换歌手、歌名或专辑关键词再试。", prompt),
			RequestTopic: prompt,
		}, nil
	}

	systemPrompt := s.buildSystemPrompt(prompt, intent, candidates, userID)
	messages := make([]AIChatMessage, 0, len(history)+1)
	messages = append(messages, AIChatMessage{Role: "system", Content: systemPrompt})

	for _, item := range history {
		role := strings.TrimSpace(item.Role)
		content := strings.TrimSpace(item.Content)
		if content == "" {
			continue
		}
		if role != "user" && role != "assistant" {
			continue
		}
		messages = append(messages, AIChatMessage{Role: role, Content: content})
	}

	return &AIChatContext{
		Messages:     messages,
		Candidates:   candidates,
		RequestTopic: prompt,
	}, nil
}

func latestUserPrompt(history []AIChatMessage) string {
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" && strings.TrimSpace(history[i].Content) != "" {
			return strings.TrimSpace(history[i].Content)
		}
	}
	return ""
}

func conversationPrompt(history []AIChatMessage) string {
	parts := make([]string, 0, 4)
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role != "user" {
			continue
		}
		content := strings.TrimSpace(history[i].Content)
		if content == "" {
			continue
		}
		parts = append([]string{content}, parts...)
		if len(parts) >= 4 {
			break
		}
	}

	if len(parts) == 0 {
		return latestUserPrompt(history)
	}
	return strings.Join(parts, " ")
}

func (s *AIService) shouldRecommendFromHistory(history []AIChatMessage) bool {
	_, intent := s.resolveRecommendationIntent(history)
	return intent.IsRecommendation
}

func (s *AIService) analyzeUserIntentAndKeywords(prompt string) *AIIntentAnalysis {
	result := &AIIntentAnalysis{Keywords: make([]string, 0)}
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return result
	}

	result.IsFollowup = containsAnyPhrase(prompt, []string{
		"再来", "还有吗", "继续", "换几首", "更多", "再推荐", "下一首",
	})
	result.Language = detectLanguageKeyword(prompt)
	result.Mood = detectMoodKeyword(prompt)
	result.Keywords = extractIntentKeywords(prompt, result.Mood, result.Language)
	result.IsRecommendation = result.IsFollowup ||
		result.Mood != "" ||
		result.Language != "" ||
		len(result.Keywords) > 0 ||
		containsAnyPhrase(prompt, []string{
			"推荐", "来几首", "来点", "想听", "适合", "有没有", "求", "听什么", "音乐", "歌曲", "歌单", "playlist", "music",
		})
	return result
}

func (s *AIService) resolveRecommendationIntent(history []AIChatMessage) (string, *AIIntentAnalysis) {
	latest := latestUserPrompt(history)
	if latest == "" {
		return "", &AIIntentAnalysis{Keywords: make([]string, 0)}
	}

	intent := s.analyzeUserIntentAndKeywords(latest)

	if intent.IsFollowup {
		if previous, prevIntent := s.findPreviousSpecificIntent(history); previous != "" {
			return previous, prevIntent
		}
	}

	if intent.IsRecommendation {
		return latest, intent
	}

	return "", intent
}

func (s *AIService) findPreviousSpecificIntent(history []AIChatMessage) (string, *AIIntentAnalysis) {
	skippedLatest := false
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role != "user" {
			continue
		}

		content := strings.TrimSpace(history[i].Content)
		if content == "" {
			continue
		}

		if !skippedLatest {
			skippedLatest = true
			continue
		}

		intent := s.analyzeUserIntentAndKeywords(content)
		if intent.IsRecommendation && !intent.IsFollowup {
			return content, intent
		}
	}

	return "", &AIIntentAnalysis{Keywords: make([]string, 0)}
}

func (s *AIService) buildSystemPrompt(prompt string, intent *AIIntentAnalysis, candidates []model.Song, userID uint) string {
	var builder strings.Builder
	builder.WriteString("你是 MusicBox 的 AI 音乐助手。请始终使用简洁、自然的中文回复。")
	builder.WriteString("当用户询问推荐歌曲时，只能推荐当前数据库候选歌曲里真实存在的内容，不要编造库里没有的歌。")
	builder.WriteString("回答推荐请求时，只允许从候选歌曲中挑选，并附上一句简短推荐理由。候选不够就少推荐，不要补充候选之外的歌曲。")
	builder.WriteString("对于伤感、治愈、纯音乐、适合写代码等语义型请求，可以根据候选歌曲的歌名、歌手、专辑信息做合理推断，但不要把明显不符合的歌说成符合。")

	if intent.IsRecommendation && len(candidates) == 0 {
		builder.WriteString("如果当前数据库候选为空，说明本轮请求没有直接命中的歌曲。不要根据历史偏好硬凑推荐，更不要把不相关歌曲包装成符合当前主题的结果。")
		builder.WriteString("此时应直接说明当前曲库缺少与本轮请求直接匹配的歌曲，并建议用户换歌手、歌名或专辑关键词。")
	}

	if intent.Language != "" {
		builder.WriteString(fmt.Sprintf("用户明确要求 `%s` 歌曲。如果候选歌曲列表为空或不符合语言要求，请告知用户库内暂无该语言歌曲，不要推荐其他语言。\n", intent.Language))
	}
	if intent.Mood != "" {
		builder.WriteString(fmt.Sprintf("当前用户寻找的心情/风格是：`%s`。请优先推荐语义符合该心情的歌曲。\n", intent.Mood))
	}

	if userID > 0 && len(candidates) > 0 {
		builder.WriteString("如果下方有用户偏好上下文，请结合最近播放和点赞习惯排序推荐结果。")
	}

	songContext := formatSongContext(candidates)
	if songContext != "" {
		builder.WriteString("\n\n当前数据库候选歌曲：\n")
		builder.WriteString(songContext)
	}

	preferenceContext := ""
	if len(candidates) > 0 {
		preferenceContext = s.buildUserPreferenceContext(userID)
	}
	if preferenceContext != "" {
		builder.WriteString("\n\n用户偏好上下文：\n")
		builder.WriteString(preferenceContext)
	}

	return builder.String()
}

func (s *AIService) buildUserPreferenceContext(userID uint) string {
	if userID == 0 {
		return ""
	}

	parts := make([]string, 0, 2)

	recentSongs, err := s.getRecentHistorySongs(userID, 5)
	if err == nil && len(recentSongs) > 0 {
		lines := make([]string, 0, len(recentSongs))
		for _, song := range recentSongs {
			lines = append(lines, fmt.Sprintf("- 最近播放: %s - %s", song.Title, song.Artist))
		}
		parts = append(parts, strings.Join(lines, "\n"))
	}

	likedSongs, err := (&SongLikesService{}).GetLikedSongsByUser(userID)
	if err == nil && len(likedSongs) > 0 {
		if len(likedSongs) > 5 {
			likedSongs = likedSongs[:5]
		}
		lines := make([]string, 0, len(likedSongs))
		for _, song := range likedSongs {
			lines = append(lines, fmt.Sprintf("- 已点赞: %s - %s", song.Title, song.Artist))
		}
		parts = append(parts, strings.Join(lines, "\n"))
	}

	return strings.Join(parts, "\n")
}

func (s *AIService) getRecommendationCandidates(prompt string, intent *AIIntentAnalysis, userID uint, limit int) ([]model.Song, error) {
	if limit <= 0 {
		limit = 4
	}

	// 尝试关键词严格匹配
	strictSongs, err := s.getStrictMatchedSongs(prompt, intent.Keywords, userID, limit)
	if err == nil && len(strictSongs) > 0 {
		// 如果关键词明确命中了（例如搜歌手、歌名），直接返回
		return strictSongs, nil
	}

	// 如果是语义型请求（有情感/风格/语言要求），获取更多潜在候选进行 AI 筛选
	fetchLimit := maxInt(limit*4, 20)
	candidates, err := s.getSemanticCandidates(userID, fetchLimit)
	if err != nil || len(candidates) == 0 {
		return nil, nil
	}

	// 轻量过滤，避免在正式流式回复前再次额外调用大模型。
	if intent.Language != "" {
		filtered := filterCandidatesByLanguage(candidates, intent.Language)
		if len(filtered) > 0 {
			candidates = filtered
		}
	}

	// 最后兜底：返回样本的前几个
	return candidates[:minInt(len(candidates), limit)], nil
}

func containsAnyPhrase(text string, phrases []string) bool {
	lower := strings.ToLower(strings.TrimSpace(text))
	for _, phrase := range phrases {
		if phrase == "" {
			continue
		}
		if strings.Contains(lower, strings.ToLower(phrase)) {
			return true
		}
	}
	return false
}

func detectLanguageKeyword(prompt string) string {
	switch {
	case containsAnyPhrase(prompt, []string{"中文", "国语", "华语"}):
		return "中文"
	case containsAnyPhrase(prompt, []string{"英文", "英语", "english"}):
		return "英文"
	case containsAnyPhrase(prompt, []string{"日文", "日语", "japanese"}):
		return "日语"
	case containsAnyPhrase(prompt, []string{"韩文", "韩语", "korean"}):
		return "韩语"
	case containsAnyPhrase(prompt, []string{"粤语", "粤语歌"}):
		return "粤语"
	default:
		return ""
	}
}

func detectMoodKeyword(prompt string) string {
	moods := []string{
		"学习", "写作业", "自习", "专注", "放松", "安静", "治愈", "伤感", "难过",
		"emo", "开心", "快乐", "失恋", "睡前", "通勤", "开车", "跑步", "运动",
	}
	for _, mood := range moods {
		if strings.Contains(strings.ToLower(prompt), strings.ToLower(mood)) {
			return mood
		}
	}
	return ""
}

func extractIntentKeywords(prompt, mood, language string) []string {
	replacer := strings.NewReplacer(
		"推荐几首", " ",
		"推荐一下", " ",
		"推荐", " ",
		"来几首", " ",
		"来点", " ",
		"想听", " ",
		"适合", " ",
		"听的", " ",
		"音乐", " ",
		"歌曲", " ",
		"歌单", " ",
		"的歌", " ",
		"来首", " ",
		"有没有", " ",
		"求", " ",
		"安慰", " ",
		"什么", " ",
		"一首", " ",
		"几首", " ",
		"给我", " ",
		"帮我", " ",
		"播放", " ",
	)
	cleaned := replacer.Replace(prompt)
	for _, token := range []string{mood, language, "中文", "英文", "日语", "韩语", "粤语"} {
		if token != "" {
			cleaned = strings.ReplaceAll(cleaned, token, " ")
		}
	}
	parts := strings.FieldsFunc(cleaned, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r) || strings.ContainsRune("，。！？、；：,.!?/|", r)
	})
	result := make([]string, 0, len(parts))
	seen := make(map[string]struct{})
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if utfLen(part) < 2 {
			continue
		}
		if _, ok := seen[part]; ok {
			continue
		}
		seen[part] = struct{}{}
		result = append(result, part)
	}
	return result
}

func filterCandidatesByLanguage(candidates []model.Song, language string) []model.Song {
	if len(candidates) == 0 || language == "" {
		return candidates
	}
	filtered := make([]model.Song, 0, len(candidates))
	for _, song := range candidates {
		switch language {
		case "英文":
			if looksEnglishSong(song) {
				filtered = append(filtered, song)
			}
		case "中文", "粤语":
			if !looksEnglishSong(song) {
				filtered = append(filtered, song)
			}
		default:
			filtered = append(filtered, song)
		}
	}
	return filtered
}

func (s *AIService) rankSongsBySemanticWithAI(prompt string, intent *AIIntentAnalysis, songs []model.Song, limit int) ([]model.Song, error) {
	if len(songs) == 0 {
		return nil, nil
	}

	messages := []AIChatMessage{
		{
			Role: "system",
			Content: `你是一个音乐推荐筛选器。请根据用户的需求，从给出的候选歌曲列表中选出最符合要求的歌曲。
要求：
- 严格符合用户要求的情绪/风格（mood: ` + intent.Mood + `）和语言（language: ` + intent.Language + `）。
- 如果候选歌曲中没有一个符合要求，请返回空数组。
- 不要推荐明显不符合要求的歌曲（例如用户要忧郁的，不要推欢快的；用户要中文的，不要推日语）。
- 只返回 JSON 格式：{"songIds": [id1, id2...]}`,
		},
		{
			Role: "user",
			Content: fmt.Sprintf("用户请求：%s\n\n候选列表：\n%s\n\n请返回最多 %d 个最符合要求的歌曲 ID。",
				prompt, formatSongSelectionContext(songs), limit),
		},
	}

	content, err := s.requestChatCompletion(messages, false)
	if err != nil {
		return nil, err
	}

	parsed, err := parseAISongSelection(content)
	if err != nil || len(parsed.SongIDs) == 0 {
		return nil, nil
	}

	songMap := make(map[uint]model.Song)
	for _, song := range songs {
		songMap[song.ID] = song
	}

	result := make([]model.Song, 0, len(parsed.SongIDs))
	for _, id := range parsed.SongIDs {
		if song, ok := songMap[id]; ok {
			result = append(result, song)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func formatSongContext(songs []model.Song) string {
	if len(songs) == 0 {
		return "当前数据库没有可用歌曲候选，请避免编造具体歌名。"
	}

	lines := make([]string, 0, len(songs))
	for i, song := range songs {
		line := fmt.Sprintf("%d. %s - %s", i+1, song.Title, song.Artist)
		if strings.TrimSpace(song.Album) != "" {
			line += fmt.Sprintf(" | 专辑：%s", song.Album)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (s *AIService) getSemanticCandidates(userID uint, limit int) ([]model.Song, error) {
	if limit <= 0 {
		limit = 12
	}

	unique := make(map[uint]model.Song)
	artistSeen := make(map[string]int)
	order := make([]uint, 0, limit)
	addSongs := func(songs []model.Song, maxPerArtist int) {
		for _, song := range songs {
			if _, exists := unique[song.ID]; exists {
				continue
			}

			artistKey := strings.ToLower(strings.TrimSpace(song.Artist))
			if maxPerArtist > 0 && artistKey != "" && artistSeen[artistKey] >= maxPerArtist {
				continue
			}

			unique[song.ID] = song
			order = append(order, song.ID)
			if artistKey != "" {
				artistSeen[artistKey]++
			}
			if len(order) >= limit {
				return
			}
		}
	}

	if userID > 0 {
		recentSongs, err := s.getRecentHistorySongs(userID, limit*2)
		if err != nil {
			return nil, err
		}
		addSongs(recentSongs, 2)

		likedSongs, err := (&SongLikesService{}).GetLikedSongsByUser(userID)
		if err != nil {
			return nil, err
		}
		addSongs(likedSongs, 2)
	}

	if len(order) < limit {
		latestSongs, err := s.getLatestSongs(limit * 2)
		if err != nil {
			return nil, err
		}
		addSongs(latestSongs, 2)
	}

	result := make([]model.Song, 0, len(order))
	for _, id := range order {
		result = append(result, unique[id])
	}
	return result, nil
}

func (s *AIService) getStrictMatchedSongs(prompt string, keywords []string, userID uint, limit int) ([]model.Song, error) {
	if limit <= 0 {
		limit = 4
	}

	songs, err := s.searchSongsByKeywords(keywords, maxInt(limit*4, 16))
	if err != nil {
		return nil, err
	}
	if len(songs) == 0 && isEnglishRequest(prompt) {
		return s.getEnglishSongs(limit)
	}
	if len(songs) == 0 {
		return nil, nil
	}

	profile, err := s.buildUserPreferenceProfile(userID)
	if err != nil {
		return nil, err
	}

	// 之前在上方获取到的关键词
	scored := make([]scoredSongCandidate, 0, len(songs))
	bestScore := 0
	for _, song := range songs {
		score := scoreSongAgainstPrompt(song, keywords)
		score += scoreSongByPreference(song, profile)
		if score <= 0 {
			continue
		}
		if score > bestScore {
			bestScore = score
		}
		scored = append(scored, scoredSongCandidate{
			Song:  song,
			Score: score,
		})
	}
	if len(scored) == 0 {
		return nil, nil
	}

	return pickTopScoredSongs(scored, limit, maxInt(3, bestScore/2)), nil
}

func formatSongSelectionContext(songs []model.Song) string {
	lines := make([]string, 0, len(songs))
	for _, song := range songs {
		line := fmt.Sprintf("- ID=%d | %s - %s", song.ID, song.Title, song.Artist)
		if strings.TrimSpace(song.Album) != "" {
			line += fmt.Sprintf(" | 专辑：%s", song.Album)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func parseAISongSelection(content string) (*aiSongSelectionResult, error) {
	cleaned := strings.TrimSpace(content)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	start := strings.Index(cleaned, "{")
	end := strings.LastIndex(cleaned, "}")
	if start >= 0 && end > start {
		cleaned = cleaned[start : end+1]
	}

	var result aiSongSelectionResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func isEnglishRequest(prompt string) bool {
	return strings.Contains(prompt, "英文") || strings.Contains(prompt, "英语") || strings.Contains(strings.ToLower(prompt), "english")
}

func (s *AIService) getEnglishSongs(limit int) ([]model.Song, error) {
	if limit <= 0 {
		limit = 4
	}

	var songs []model.Song
	if err := global.DB.Model(&model.Song{}).
		Order("upload_at DESC").
		Limit(80).
		Find(&songs).Error; err != nil {
		return nil, err
	}

	filtered := make([]model.Song, 0, limit)
	for _, song := range songs {
		if !looksEnglishSong(song) {
			continue
		}
		filtered = append(filtered, song)
		if len(filtered) >= limit {
			break
		}
	}
	return filtered, nil
}

func looksEnglishSong(song model.Song) bool {
	return isMostlyASCII(song.Title) || isMostlyASCII(song.Artist)
}

func isMostlyASCII(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}

	total := 0
	ascii := 0
	for _, r := range value {
		if unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsDigit(r) {
			continue
		}
		total++
		if r <= unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r)) {
			ascii++
		}
	}
	if total == 0 {
		return false
	}
	return ascii*2 >= total
}

func (s *AIService) searchSongsByKeywords(keywords []string, limit int) ([]model.Song, error) {
	if len(keywords) == 0 {
		return nil, nil
	}
	if len(keywords) > 5 {
		keywords = keywords[:5]
	}

	query := global.DB.Model(&model.Song{})
	for index, keyword := range keywords {
		pattern := "%" + keyword + "%"
		clause := "title LIKE ? OR artist LIKE ? OR album LIKE ?"
		if index == 0 {
			query = query.Where(clause, pattern, pattern, pattern)
			continue
		}
		query = query.Or(clause, pattern, pattern, pattern)
	}

	var songs []model.Song
	err := query.Order("upload_at DESC").Limit(limit).Find(&songs).Error
	return songs, err
}

func scoreSongAgainstPrompt(song model.Song, keywords []string) int {
	if len(keywords) == 0 {
		return 0
	}

	title := strings.ToLower(strings.TrimSpace(song.Title))
	artist := strings.ToLower(strings.TrimSpace(song.Artist))
	album := strings.ToLower(strings.TrimSpace(song.Album))
	score := 0

	for _, keyword := range keywords {
		kw := strings.ToLower(strings.TrimSpace(keyword))
		if kw == "" {
			continue
		}

		switch {
		case artist == kw:
			score += 18
		case strings.HasPrefix(artist, kw):
			score += 14
		case strings.Contains(artist, kw):
			score += 10
		}

		switch {
		case title == kw:
			score += 15
		case strings.HasPrefix(title, kw):
			score += 11
		case strings.Contains(title, kw):
			score += 8
		}

		switch {
		case album == kw:
			score += 8
		case album != "" && strings.HasPrefix(album, kw):
			score += 5
		case album != "" && strings.Contains(album, kw):
			score += 3
		}
	}

	return score
}

func utfLen(value string) int {
	return len([]rune(value))
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func scoreSongByPreference(song model.Song, profile *userPreferenceProfile) int {
	if profile == nil {
		return 0
	}

	score := 0
	score += profile.SongWeights[song.ID]
	score += profile.ArtistWeights[strings.ToLower(strings.TrimSpace(song.Artist))]
	score += profile.AlbumWeights[strings.ToLower(strings.TrimSpace(song.Album))]
	return score
}

func (s *AIService) buildUserPreferenceProfile(userID uint) (*userPreferenceProfile, error) {
	if userID == 0 {
		return nil, nil
	}

	profile := &userPreferenceProfile{
		SongWeights:   make(map[uint]int),
		ArtistWeights: make(map[string]int),
		AlbumWeights:  make(map[string]int),
	}

	recentSongs, err := s.getRecentHistorySongs(userID, 24)
	if err != nil {
		return nil, err
	}
	for i, song := range recentSongs {
		recencyWeight := maxInt(1, 8-(i/3))
		profile.SongWeights[song.ID] += recencyWeight * 4
		addStringWeight(profile.ArtistWeights, song.Artist, recencyWeight*3)
		addStringWeight(profile.AlbumWeights, song.Album, recencyWeight)
	}

	likedSongs, err := (&SongLikesService{}).GetLikedSongsByUser(userID)
	if err != nil {
		return nil, err
	}
	for _, song := range likedSongs {
		profile.SongWeights[song.ID] += 18
		addStringWeight(profile.ArtistWeights, song.Artist, 10)
		addStringWeight(profile.AlbumWeights, song.Album, 4)
	}

	return profile, nil
}

func addStringWeight(target map[string]int, key string, weight int) {
	normalized := strings.ToLower(strings.TrimSpace(key))
	if normalized == "" || weight <= 0 {
		return
	}
	target[normalized] += weight
}

func pickTopScoredSongs(scored []scoredSongCandidate, limit int, threshold int) []model.Song {
	if len(scored) == 0 || limit <= 0 {
		return nil
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].Score == scored[j].Score {
			if scored[i].Song.UploadAt.Equal(scored[j].Song.UploadAt) {
				return scored[i].Song.ID > scored[j].Song.ID
			}
			return scored[i].Song.UploadAt.After(scored[j].Song.UploadAt)
		}
		return scored[i].Score > scored[j].Score
	})

	result := make([]model.Song, 0, minInt(limit, len(scored)))
	for _, item := range scored {
		if threshold > 0 && item.Score < threshold {
			continue
		}
		result = append(result, item.Song)
		if len(result) >= limit {
			break
		}
	}

	if len(result) == 0 && len(scored) > 0 {
		for _, item := range scored[:minInt(limit, len(scored))] {
			result = append(result, item.Song)
			if len(result) >= limit {
				break
			}
		}
	}

	return result
}

func (s *AIService) getRecentHistorySongs(userID uint, limit int) ([]model.Song, error) {
	var songs []model.Song
	err := global.DB.
		Table("songs").
		Joins("JOIN play_histories ON play_histories.song_id = songs.id").
		Where("play_histories.user_id = ?", userID).
		Order("play_histories.played_at DESC").
		Limit(limit).
		Find(&songs).Error
	return songs, err
}

func (s *AIService) getLatestSongs(limit int) ([]model.Song, error) {
	var songs []model.Song
	err := global.DB.Model(&model.Song{}).Order("upload_at DESC").Limit(limit).Find(&songs).Error
	return songs, err
}
