package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// CloudflareWorkersAIProvider Cloudflare Workers AI Provider
// 通过 AI Gateway（推荐）或直接 REST API 调用 Workers AI 模型
type CloudflareWorkersAIProvider struct {
	enabled    bool
	accountID  string
	apiToken   string
	gatewayID  string
	useGateway bool
	model      string
	client     *http.Client
}

// NewCloudflareWorkersAIProvider 创建 Cloudflare Workers AI Provider
func NewCloudflareWorkersAIProvider(enabled bool, accountID, apiToken, gatewayID string, useGateway bool, model string, timeout time.Duration) *CloudflareWorkersAIProvider {
	return &CloudflareWorkersAIProvider{
		enabled:    enabled,
		accountID:  accountID,
		apiToken:   apiToken,
		gatewayID:  gatewayID,
		useGateway: useGateway,
		model:      model,
		client:     &http.Client{Timeout: timeout},
	}
}

func (p *CloudflareWorkersAIProvider) Name() ProviderName { return ProviderCloudflareWorkersAI }

func (p *CloudflareWorkersAIProvider) Enabled() bool {
	if !p.enabled || p.accountID == "" || p.apiToken == "" || p.model == "" {
		return false
	}
	if p.useGateway && p.gatewayID == "" {
		return false
	}
	return true
}

func (p *CloudflareWorkersAIProvider) endpoint() string {
	escapedModel := url.PathEscape(p.model)
	if p.useGateway {
		return fmt.Sprintf(
			"https://gateway.ai.cloudflare.com/v1/%s/%s/workers-ai/%s",
			p.accountID,
			p.gatewayID,
			escapedModel,
		)
	}
	return fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s",
		p.accountID,
		escapedModel,
	)
}

func (p *CloudflareWorkersAIProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	payload := map[string]any{
		"messages": req.Messages,
		"stream":   false,
	}
	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}
	if req.Temperature > 0 {
		payload["temperature"] = req.Temperature
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+p.apiToken)
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
		return nil, fmt.Errorf("cloudflare workers ai error: status=%d body=%s", resp.StatusCode, string(raw))
	}

	var out struct {
		Success bool `json:"success"`
		Result  struct {
			Response string `json:"response"`
		} `json:"result"`
		Usage struct {
			PromptTokens        int `json:"prompt_tokens"`
			CompletionTokens    int `json:"completion_tokens"`
			TotalTokens         int `json:"total_tokens"`
			PromptTokensDetails struct {
				CachedTokens int `json:"cached_tokens"`
			} `json:"prompt_tokens_details"`
		} `json:"usage"`
		Errors []any `json:"errors"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.Result.Response == "" {
		return nil, fmt.Errorf("cloudflare workers ai empty response: %s", string(raw))
	}

	usage := EstimateKimiNeurons(req.Messages, out.Result.Response)
	if out.Usage.TotalTokens > 0 || out.Usage.PromptTokens > 0 || out.Usage.CompletionTokens > 0 {
		usage.InputTokens = out.Usage.PromptTokens
		usage.OutputTokens = out.Usage.CompletionTokens
		usage.CachedTokens = out.Usage.PromptTokensDetails.CachedTokens
		usage.TotalTokens = out.Usage.TotalTokens
		if usage.TotalTokens == 0 {
			usage.TotalTokens = usage.InputTokens + usage.OutputTokens
		}
	}

	return &ChatResponse{
		Content:  out.Result.Response,
		Provider: ProviderCloudflareWorkersAI,
		Model:    p.model,
		Usage:    usage,
	}, nil
}
