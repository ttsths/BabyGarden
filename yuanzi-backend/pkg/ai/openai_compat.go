package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenAICompatProvider 通用 OpenAI-compatible API Provider
// 适用于 GrokAI / CLIProxyAPI 等兼容 OpenAI chat/completions 接口的服务
type OpenAICompatProvider struct {
	name    ProviderName
	enabled bool
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

// NewOpenAICompatProvider 创建 OpenAI 兼容 Provider
func NewOpenAICompatProvider(name ProviderName, enabled bool, baseURL, apiKey, model string, timeout time.Duration) *OpenAICompatProvider {
	return &OpenAICompatProvider{
		name:    name,
		enabled: enabled,
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		model:   model,
		client:  &http.Client{Timeout: timeout},
	}
}

func (p *OpenAICompatProvider) Name() ProviderName { return p.name }

func (p *OpenAICompatProvider) Enabled() bool {
	return p.enabled && p.baseURL != "" && p.apiKey != ""
}

func (p *OpenAICompatProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	payload := map[string]any{
		"model":       p.model,
		"messages":    req.Messages,
		"stream":      false,
		"temperature": req.Temperature,
	}
	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	if req.RequestID != "" {
		httpReq.Header.Set("X-Request-ID", req.RequestID)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s error: status=%d body=%s", p.name, resp.StatusCode, string(raw))
	}

	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens        int `json:"prompt_tokens"`
			CompletionTokens    int `json:"completion_tokens"`
			TotalTokens         int `json:"total_tokens"`
			PromptTokensDetails struct {
				CachedTokens int `json:"cached_tokens"`
			} `json:"prompt_tokens_details"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if len(out.Choices) == 0 {
		return nil, fmt.Errorf("%s empty choices", p.name)
	}

	return &ChatResponse{
		Content:  out.Choices[0].Message.Content,
		Provider: p.name,
		Model:    p.model,
		Usage: Usage{
			InputTokens:  out.Usage.PromptTokens,
			OutputTokens: out.Usage.CompletionTokens,
			CachedTokens: out.Usage.PromptTokensDetails.CachedTokens,
			TotalTokens:  out.Usage.TotalTokens,
		},
	}, nil
}
