package ai

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestEstimateTokenCount(t *testing.T) {
	tests := []struct {
		input    string
		minExpected int
	}{
		{"", 0},
		{"hello", 1},
		{"你好世界", 1},
		{"hello你好", 1},
	}

	for _, tt := range tests {
		got := CountTokens(tt.input, "qwen-plus")
		if got <= 0 && tt.input != "" {
			t.Errorf("CountTokens(%q) = %d, want > 0", tt.input, got)
		}
		if tt.input == "" && got != 0 {
			t.Errorf("CountTokens(empty) = %d, want 0", got)
		}
	}
}

func TestEstimateKimiNeurons(t *testing.T) {
	messages := []ChatMessage{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: "Hello, how are you?"},
	}
	output := "I'm doing well, thank you!"

	usage := EstimateKimiNeurons(messages, output)

	if usage.InputTokens <= 0 {
		t.Errorf("Expected positive input tokens, got %d", usage.InputTokens)
	}
	if usage.OutputTokens <= 0 {
		t.Errorf("Expected positive output tokens, got %d", usage.OutputTokens)
	}
	if usage.TotalTokens != usage.InputTokens+usage.OutputTokens {
		t.Errorf("TotalTokens mismatch: %d != %d + %d", usage.TotalTokens, usage.InputTokens, usage.OutputTokens)
	}
	if usage.NeuronsEst <= 0 {
		t.Errorf("Expected positive NeuronsEst, got %d", usage.NeuronsEst)
	}
}

func TestNewOpenAICompatProvider(t *testing.T) {
	p := NewOpenAICompatProvider(ProviderGrokAI, true, "https://api.example.com", "sk-test", "test-model", 30e9)
	if p.Name() != ProviderGrokAI {
		t.Errorf("Name() = %s, want %s", p.Name(), ProviderGrokAI)
	}
	if !p.Enabled() {
		t.Error("Expected provider to be enabled")
	}

	p2 := NewOpenAICompatProvider(ProviderCLIProxyAPI, false, "", "", "", 0)
	if p2.Enabled() {
		t.Error("Expected disabled provider")
	}
}

func TestOpenAICompatProviderChatUsage(t *testing.T) {
	p := NewOpenAICompatProvider(ProviderGrokAI, true, "https://api.example.test", "sk-test", "grok-test", 30e9)
	p.client = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"choices":[{"message":{"content":"ok"}}],
				"usage":{
					"prompt_tokens":10,
					"completion_tokens":4,
					"total_tokens":14,
					"prompt_tokens_details":{"cached_tokens":3}
				}
			}`)),
		}, nil
	})}

	resp, err := p.Chat(context.Background(), ChatRequest{Messages: []ChatMessage{{Role: "user", Content: "hi"}}})
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}
	if resp.Usage.InputTokens != 10 || resp.Usage.OutputTokens != 4 || resp.Usage.CachedTokens != 3 || resp.Usage.TotalTokens != 14 {
		t.Fatalf("usage解析错误: %+v", resp.Usage)
	}
}

func TestNewCloudflareWorkersAIProvider(t *testing.T) {
	p := NewCloudflareWorkersAIProvider(true, "acct-123", "token-456", "gw-789", true, "@cf/moonshotai/kimi-k2.6", 30e9)
	if p.Name() != ProviderCloudflareWorkersAI {
		t.Errorf("Name() = %s, want %s", p.Name(), ProviderCloudflareWorkersAI)
	}
	if !p.Enabled() {
		t.Error("Expected provider to be enabled")
	}

	// Missing gateway when useGateway=true
	p2 := NewCloudflareWorkersAIProvider(true, "acct-123", "token-456", "", true, "model", 0)
	if p2.Enabled() {
		t.Error("Expected disabled when gateway missing and useGateway=true")
	}

	// Disabled
	p3 := NewCloudflareWorkersAIProvider(false, "", "", "", false, "", 0)
	if p3.Enabled() {
		t.Error("Expected disabled provider")
	}
}

func TestCloudflareWorkersAIProviderChatUsage(t *testing.T) {
	p := NewCloudflareWorkersAIProvider(true, "acct-123", "token-456", "", false, "@cf/test/model", 30e9)
	p.client = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"success": true,
				"result": {"response": "ok"},
				"usage": {
					"prompt_tokens": 8,
					"completion_tokens": 5,
					"total_tokens": 13,
					"prompt_tokens_details": {"cached_tokens": 2}
				}
			}`)),
		}, nil
	})}

	resp, err := p.Chat(context.Background(), ChatRequest{Messages: []ChatMessage{{Role: "user", Content: "hi"}}})
	if err != nil {
		t.Fatalf("Chat failed: %v", err)
	}
	if resp.Usage.InputTokens != 8 || resp.Usage.OutputTokens != 5 || resp.Usage.CachedTokens != 2 || resp.Usage.TotalTokens != 13 {
		t.Fatalf("usage解析错误: %+v", resp.Usage)
	}
}

func TestDashScopeProvider(t *testing.T) {
	p := NewDashScopeProvider(true, "sk-test", "qwen-plus", 30*time.Second)
	if p.Name() != ProviderDashScope {
		t.Errorf("Name() = %s, want %s", p.Name(), ProviderDashScope)
	}
	if !p.Enabled() {
		t.Error("Expected enabled provider")
	}

	p2 := NewDashScopeProvider(true, "", "", 30*time.Second)
	if p2.Enabled() {
		t.Error("Expected disabled when API key missing")
	}
}

type mockProvider struct {
	name    ProviderName
	enabled bool
	chatFn  func(req ChatRequest) (*ChatResponse, error)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func (m *mockProvider) Name() ProviderName { return m.name }
func (m *mockProvider) Enabled() bool      { return m.enabled }
func (m *mockProvider) Chat(_ context.Context, req ChatRequest) (*ChatResponse, error) {
	return m.chatFn(req)
}

func TestRouterProviderOrder(t *testing.T) {
	callOrder := make([]ProviderName, 0)

	p1 := &mockProvider{
		name:    "grokai",
		enabled: true,
		chatFn: func(req ChatRequest) (*ChatResponse, error) {
			callOrder = append(callOrder, "grokai")
			return &ChatResponse{Content: "from grok", Provider: "grokai", Model: "grok-4"}, nil
		},
	}
	p2 := &mockProvider{
		name:    "cloudflare",
		enabled: true,
		chatFn: func(req ChatRequest) (*ChatResponse, error) {
			callOrder = append(callOrder, "cloudflare")
			return &ChatResponse{Content: "from cf", Provider: "cloudflare", Model: "kimi"}, nil
		},
	}

	router := NewRouter(nil, p1, p2)
	resp, err := router.Chat(context.TODO(), ChatRequest{})
	if err != nil {
		t.Fatalf("Router.Chat failed: %v", err)
	}
	if resp.Content != "from grok" {
		t.Errorf("Expected 'from grok', got '%s'", resp.Content)
	}
	if len(callOrder) != 1 || callOrder[0] != "grokai" {
		t.Errorf("Only grokai should be called, got %v", callOrder)
	}
}

func TestRouterFallback(t *testing.T) {
	p1 := &mockProvider{
		name:    "grokai",
		enabled: true,
		chatFn: func(req ChatRequest) (*ChatResponse, error) {
			return nil, fmt.Errorf("grokai down")
		},
	}
	p2 := &mockProvider{
		name:    "cloudflare",
		enabled: true,
		chatFn: func(req ChatRequest) (*ChatResponse, error) {
			return &ChatResponse{Content: "from cf", Provider: "cloudflare", Model: "kimi"}, nil
		},
	}

	router := NewRouter(nil, p1, p2)
	resp, err := router.Chat(context.TODO(), ChatRequest{})
	if err != nil {
		t.Fatalf("Router.Chat fallback failed: %v", err)
	}
	if resp.Provider != "cloudflare" {
		t.Errorf("Expected fallback to cloudflare, got %s", resp.Provider)
	}
}

func TestRouterAllFailed(t *testing.T) {
	p1 := &mockProvider{
		name:    "grokai",
		enabled: true,
		chatFn: func(req ChatRequest) (*ChatResponse, error) {
			return nil, fmt.Errorf("error 1")
		},
	}
	p2 := &mockProvider{
		name:    "cloudflare",
		enabled: true,
		chatFn: func(req ChatRequest) (*ChatResponse, error) {
			return nil, fmt.Errorf("error 2")
		},
	}

	router := NewRouter(nil, p1, p2)
	_, err := router.Chat(context.TODO(), ChatRequest{})
	if err == nil {
		t.Fatal("Expected error when all providers fail")
	}
}

func TestRouterSkipsDisabled(t *testing.T) {
	p1 := &mockProvider{
		name:    "grokai",
		enabled: false,
		chatFn: func(req ChatRequest) (*ChatResponse, error) {
			return nil, fmt.Errorf("should not be called")
		},
	}
	p2 := &mockProvider{
		name:    "cloudflare",
		enabled: true,
		chatFn: func(req ChatRequest) (*ChatResponse, error) {
			return &ChatResponse{Content: "from cf", Provider: "cloudflare", Model: "kimi"}, nil
		},
	}

	router := NewRouter(nil, p1, p2)
	resp, err := router.Chat(context.TODO(), ChatRequest{})
	if err != nil {
		t.Fatalf("Router.Chat failed: %v", err)
	}
	if resp.Provider != "cloudflare" {
		t.Errorf("Expected cloudflare, got %s", resp.Provider)
	}
}
