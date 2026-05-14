package ai

import (
	"fmt"
	"sync"
	"unicode/utf8"

	tiktoken "github.com/pkoukk/tiktoken-go"
)

// Kimi K2.6 Neurons 估算常量
// 基于 Cloudflare 官方定价换算：
//
//	$0.95  / M input tokens   → 86,364 Neurons / M tokens
//	$4.00  / M output tokens  → 363,636 Neurons / M tokens
//	换算基准: $0.011 / 1000 Neurons
const (
	KimiInputNeuronsPerMillion  = 86364
	KimiOutputNeuronsPerMillion = 363636
)

// tokenizerCache 缓存已初始化的 tiktoken 实例，避免重复加载编码表
var (
	tokenizerMu     sync.RWMutex
	tokenizerCache   = make(map[string]*tiktoken.Tiktoken)
	tokenizerInitErr = make(map[string]error)
)

// getTokenizer 获取或初始化指定 encoding 的 tiktoken 实例
func getTokenizer(encoding string) (*tiktoken.Tiktoken, error) {
	tokenizerMu.RLock()
	if t, ok := tokenizerCache[encoding]; ok {
		tokenizerMu.RUnlock()
		return t, nil
	}
	if err, ok := tokenizerInitErr[encoding]; ok {
		tokenizerMu.RUnlock()
		return nil, err
	}
	tokenizerMu.RUnlock()

	tokenizerMu.Lock()
	defer tokenizerMu.Unlock()

	// double-check after acquiring write lock
	if t, ok := tokenizerCache[encoding]; ok {
		return t, nil
	}

	tke, err := tiktoken.EncodingForModel(encoding)
	if err != nil {
		// 尝试按 encoding name 获取
		tke, err = tiktoken.GetEncoding(encoding)
		if err != nil {
			tokenizerInitErr[encoding] = fmt.Errorf("tiktoken init failed for %s: %w", encoding, err)
			return nil, tokenizerInitErr[encoding]
		}
	}
	tokenizerCache[encoding] = tke
	return tke, nil
}

// modelToEncoding 将 AI 模型名映射到 tiktoken encoding
// 不同模型使用不同的 tokenizer：
//   - GPT-4/GPT-3.5 系列 → cl100k_base
//   - GPT-4o/o1/o3 系列 → o200k_e
//   - qwen 系列 → cl100k_base（近似，qwen 无官方 tiktoken）
//   - grok 系列 → o200k_e（近似）
//   - kimi 系列 → cl100k_base（近似）
func modelToEncoding(model string) string {
	switch {
	case containsAny(model, "gpt-4o", "gpt-4o-", "o1", "o3", "grok-4", "grok-3"):
		return "o200k_e"
	case containsAny(model, "gpt-4", "gpt-3.5", "qwen", "kimi", "@cf/moonshotai"):
		return "cl100k_base"
	default:
		return "cl100k_base" // 默认 fallback
	}
}

func containsAny(s string, patterns ...string) bool {
	for _, p := range patterns {
		if len(s) >= len(p) && s[:len(p)] == p {
			return true
		}
	}
	return false
}

// CountTokens 使用 tiktoken 精确计算 token 数量
// 如果 tiktoken 初始化失败，fallback 到粗估
func CountTokens(text string, model string) int {
	encoding := modelToEncoding(model)
	tke, err := getTokenizer(encoding)
	if err != nil {
		return fallbackEstimateTokenCount(text)
	}
	return len(tke.Encode(text, []string{"all"}, nil))
}

// CountMessagesTokens 计算一组消息的 token 总数
func CountMessagesTokens(messages []ChatMessage, model string) int {
	total := 0
	for _, m := range messages {
		// OpenAI 格式：每条消息额外约 4 tokens overhead（role + separators）
		total += CountTokens(m.Content, model) + 4
	}
	// 消息整体约 2 tokens overhead（priming + terminator）
	total += 2
	return total
}

// fallbackEstimateTokenCount 粗估 token 数量（tiktoken 不可用时）
// 中英文混合场景下保守估算: 1 token ≈ 2 chars
func fallbackEstimateTokenCount(s string) int {
	chars := utf8.RuneCountInString(s)
	if chars == 0 {
		return 0
	}
	return chars/2 + 1
}

// EstimateKimiNeurons 按 Kimi K2.6 定价估算 Neurons 消耗
// 优先使用 tiktoken 精确计数，fallback 到粗估
func EstimateKimiNeurons(messages []ChatMessage, output string) Usage {
	model := "@cf/moonshotai/kimi-k2.6"
	inputTokens := CountMessagesTokens(messages, model)
	outputTokens := CountTokens(output, model)

	inputNeurons := inputTokens * KimiInputNeuronsPerMillion / 1_000_000
	outputNeurons := outputTokens * KimiOutputNeuronsPerMillion / 1_000_000

	return Usage{
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		TotalTokens:  inputTokens + outputTokens,
		NeuronsEst:   inputNeurons + outputNeurons,
	}
}