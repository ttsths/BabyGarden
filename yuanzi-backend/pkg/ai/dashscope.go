package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"yuanzi-backend/config"
	"yuanzi-backend/logger"
)

const (
	dashScopeAPIURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
)

// ============================================================================
// DashScopeProvider — 实现统一 Provider 接口
// ============================================================================

// DashScopeProvider 阿里云 DashScope Provider
type DashScopeProvider struct {
	enabled bool
	apiKey  string
	model   string
}

// NewDashScopeProvider 创建 DashScope Provider
func NewDashScopeProvider(enabled bool, apiKey, model string) *DashScopeProvider {
	return &DashScopeProvider{enabled: enabled, apiKey: apiKey, model: model}
}

func (p *DashScopeProvider) Name() ProviderName { return ProviderDashScope }

func (p *DashScopeProvider) Enabled() bool {
	return p.enabled && p.apiKey != ""
}

func (p *DashScopeProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	dashMessages := make([]ChatMessage, len(req.Messages))
	for i, m := range req.Messages {
		dashMessages[i] = ChatMessage{Role: m.Role, Content: m.Content}
	}

	client := &dashClient{apiKey: p.apiKey}
	resp, err := client.chat(dashMessages, p.model, req.MaxTokens, req.Temperature)
	if err != nil {
		return nil, err
	}

	return &ChatResponse{
		Content:  resp.Output.Text,
		Provider: ProviderDashScope,
		Model:    p.model,
		Usage: Usage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
			CachedTokens: resp.Usage.CachedTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
	}, nil
}

// ============================================================================
// 旧版 Client 与类型（保持向后兼容）
// ============================================================================

// LegacyChatResponse DashScope 响应（旧版类型，向后兼容）
type LegacyChatResponse struct {
	Output struct {
		Text string `json:"text"`
	} `json:"output"`
	Usage struct {
		InputTokens          int `json:"input_tokens"`
		OutputTokens         int `json:"output_tokens"`
		CachedTokens         int `json:"cached_tokens"`
		PromptCacheHitTokens int `json:"prompt_cache_hit_tokens"`
		TotalTokens          int `json:"total_tokens"`
	} `json:"usage"`
}

// Client AI客户端（旧版，用于向后兼容）
type Client struct {
	apiKey string
}

// NewClient 创建AI客户端（旧版）
func NewClient() *Client {
	return &Client{
		apiKey: config.GlobalConfig.AI.DashScopeAPIKey,
	}
}

// Chat 旧版 Chat（向后兼容）
func (c *Client) Chat(messages []ChatMessage) (*LegacyChatResponse, error) {
	return (&dashClient{apiKey: c.apiKey}).chat(messages, "qwen-turbo", 1000, 0.7)
}

// ============================================================================
// DashScope 内部类型
// ============================================================================

// dashScopeAPIReq DashScope API 请求体
type dashScopeAPIReq struct {
	Model      string         `json:"model"`
	Input      dashInput      `json:"input"`
	Parameters dashParameters `json:"parameters"`
}

type dashInput struct {
	Messages []ChatMessage `json:"messages"`
}

type dashParameters struct {
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// dashClient 内部客户端
type dashClient struct {
	apiKey string
}

// chat 内部实现
func (c *dashClient) chat(messages []ChatMessage, model string, maxTokens int, temperature float64) (*LegacyChatResponse, error) {
	if model == "" {
		model = "qwen-turbo"
	}
	if maxTokens <= 0 {
		maxTokens = 1000
	}
	if temperature <= 0 {
		temperature = 0.7
	}

	reqBody := dashScopeAPIReq{
		Model: model,
		Input: dashInput{
			Messages: messages,
		},
		Parameters: dashParameters{
			Temperature: temperature,
			MaxTokens:   maxTokens,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", dashScopeAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("AI chat request failed", logger.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI API error: %s", string(body))
	}

	var result LegacyChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.Usage.CachedTokens == 0 && result.Usage.PromptCacheHitTokens > 0 {
		result.Usage.CachedTokens = result.Usage.PromptCacheHitTokens
	}

	return &result, nil
}

// BuildSystemPrompt 构建系统提示（旧版）
func BuildSystemPrompt() string {
	return `你是「小园子」育儿助手，专注于0-3岁婴幼儿护理。
回答要求：
1. 简洁易懂，适合新手父母
2. 必要时给出具体操作步骤
3. 涉及医疗建议时添加免责声明
4. 不确定时建议咨询专业医生`
}
