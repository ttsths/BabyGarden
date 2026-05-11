package ai

import "context"

// ProviderName 标识 AI 服务提供商
type ProviderName string

const (
	ProviderGrokAI              ProviderName = "grokai"
	ProviderCloudflareWorkersAI ProviderName = "cloudflare_workers_ai"
	ProviderDashScope           ProviderName = "dashscope"
	ProviderCLIProxyAPI         ProviderName = "cliproxyapi"
)

// ChatMessage 统一聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 统一聊天请求
type ChatRequest struct {
	UserID      uint          `json:"user_id"`
	FamilyID    uint          `json:"family_id"`
	Messages    []ChatMessage `json:"messages"`
	Stream      bool          `json:"stream"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	RequestID   string        `json:"request_id,omitempty"`
}

// Usage token 与额度用量
type Usage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
	TotalTokens  int `json:"total_tokens,omitempty"`
	NeuronsEst   int `json:"neurons_est,omitempty"`
}

// ChatResponse 统一聊天响应
type ChatResponse struct {
	Content  string       `json:"content"`
	Provider ProviderName `json:"provider"`
	Model    string       `json:"model"`
	Usage    Usage        `json:"usage,omitempty"`
}

// Provider 统一 AI Provider 接口
type Provider interface {
	Name() ProviderName
	Enabled() bool
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}
