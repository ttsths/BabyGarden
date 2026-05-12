package ai

import (
	"unicode/utf8"
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

// EstimateTokenCount 粗估 token 数量
// 中英文混合场景下保守估算: 1 token ≈ 2 chars
// 仅用于额度闸门，不用于计费结算
func EstimateTokenCount(s string) int {
	chars := utf8.RuneCountInString(s)
	if chars == 0 {
		return 0
	}
	return chars/2 + 1
}

// EstimateKimiNeurons 按 Kimi K2.6 定价估算 Neurons 消耗
func EstimateKimiNeurons(messages []ChatMessage, output string) Usage {
	inputChars := 0
	for _, m := range messages {
		inputChars += utf8.RuneCountInString(m.Content)
	}
	inputTokens := inputChars/2 + 1
	outputTokens := EstimateTokenCount(output)

	inputNeurons := inputTokens * KimiInputNeuronsPerMillion / 1_000_000
	outputNeurons := outputTokens * KimiOutputNeuronsPerMillion / 1_000_000

	return Usage{
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		TotalTokens:  inputTokens + outputTokens,
		NeuronsEst:   inputNeurons + outputNeurons,
	}
}
