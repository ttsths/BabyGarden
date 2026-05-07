package ai

import (
	"bytes"
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

// Client AI客户端
type Client struct {
	apiKey string
}

// NewClient 创建AI客户端
func NewClient() *Client {
	return &Client{
		apiKey: config.GlobalConfig.AI.DashScopeAPIKey,
	}
}

// ChatMessage 对话消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 对话请求
type ChatRequest struct {
	Model      string        `json:"model"`
	Input      ChatInput     `json:"input"`
	Parameters ChatParameters `json:"parameters"`
}

// ChatInput 输入
type ChatInput struct {
	Messages []ChatMessage `json:"messages"`
}

// ChatParameters 参数
type ChatParameters struct {
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// ChatResponse 对话响应
type ChatResponse struct {
	Output struct {
		Text string `json:"text"`
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// Chat 对话
func (c *Client) Chat(messages []ChatMessage) (*ChatResponse, error) {
	reqBody := ChatRequest{
		Model: "qwen-turbo",
		Input: ChatInput{
			Messages: messages,
		},
		Parameters: ChatParameters{
			Temperature: 0.7,
			MaxTokens:   1000,
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

	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// BuildSystemPrompt 构建系统提示
func BuildSystemPrompt() string {
	return `你是「小园子」育儿助手，专注于0-3岁婴幼儿护理。
回答要求：
1. 简洁易懂，适合新手父母
2. 必要时给出具体操作步骤
3. 涉及医疗建议时添加免责声明
4. 不确定时建议咨询专业医生`
}
