// Package tokencount provides token counting and cost estimation for LLM APIs.
// It implements a simple BPE-approximate tokenizer and maintains a database of model pricing.
package tokencount

import (
	"strings"
	"unicode"
)

// ModelInfo contains pricing and context window information for an LLM model.
type ModelInfo struct {
	Name            string  `json:"name"`
	Provider        string  `json:"provider"`
	MaxTokens       int     `json:"max_tokens"`
	InputCostPer1M  float64 `json:"input_cost_per_1m"`
	OutputCostPer1M float64 `json:"output_cost_per_1m"`
}

// CostEstimate is the result of a cost estimation.
type CostEstimate struct {
	Model        string  `json:"model"`
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	InputCost    float64 `json:"input_cost"`
	OutputCost   float64 `json:"output_cost"`
	TotalCost    float64 `json:"total_cost"`
}

var models = map[string]ModelInfo{
	"gpt-4o":            {Name: "gpt-4o", Provider: "openai", MaxTokens: 128000, InputCostPer1M: 2.5, OutputCostPer1M: 10},
	"gpt-4o-mini":       {Name: "gpt-4o-mini", Provider: "openai", MaxTokens: 128000, InputCostPer1M: 0.15, OutputCostPer1M: 0.6},
	"gpt-4-turbo":       {Name: "gpt-4-turbo", Provider: "openai", MaxTokens: 128000, InputCostPer1M: 10, OutputCostPer1M: 30},
	"gpt-4":             {Name: "gpt-4", Provider: "openai", MaxTokens: 8192, InputCostPer1M: 30, OutputCostPer1M: 60},
	"gpt-3.5-turbo":     {Name: "gpt-3.5-turbo", Provider: "openai", MaxTokens: 16385, InputCostPer1M: 0.5, OutputCostPer1M: 1.5},
	"claude-3.5-sonnet": {Name: "claude-3.5-sonnet", Provider: "anthropic", MaxTokens: 200000, InputCostPer1M: 3, OutputCostPer1M: 15},
	"claude-3-opus":     {Name: "claude-3-opus", Provider: "anthropic", MaxTokens: 200000, InputCostPer1M: 15, OutputCostPer1M: 75},
	"claude-3-haiku":    {Name: "claude-3-haiku", Provider: "anthropic", MaxTokens: 200000, InputCostPer1M: 0.25, OutputCostPer1M: 1.25},
	"gemini-1.5-pro":    {Name: "gemini-1.5-pro", Provider: "google", MaxTokens: 2097152, InputCostPer1M: 1.25, OutputCostPer1M: 5},
	"gemini-2.0-flash":  {Name: "gemini-2.0-flash", Provider: "google", MaxTokens: 1048576, InputCostPer1M: 0.1, OutputCostPer1M: 0.4},
}

// CountTokens returns an approximate token count for the given text.
// Uses a simple word/subword splitting heuristic that approximates BPE tokenization.
func CountTokens(text string) int {
	if text == "" {
		return 0
	}
	count := 0
	words := strings.Fields(text)
	for _, word := range words {
		// Approximate: short words = 1 token, longer words split ~4 chars per token
		runes := []rune(word)
		wLen := len(runes)
		if wLen <= 4 {
			count++
		} else {
			count += (wLen + 3) / 4
		}
		// Punctuation attached to words counts extra
		for _, r := range runes {
			if unicode.IsPunct(r) || unicode.IsSymbol(r) {
				count++
			}
		}
	}
	return count
}

// GetModel returns model information by name.
func GetModel(name string) (ModelInfo, bool) {
	m, ok := models[name]
	return m, ok
}

// ListModels returns all known models.
func ListModels() []ModelInfo {
	result := make([]ModelInfo, 0, len(models))
	for _, m := range models {
		result = append(result, m)
	}
	return result
}

// EstimateCost calculates the cost for a given model, input tokens, and output tokens.
func EstimateCost(modelName string, inputTokens, outputTokens int) (CostEstimate, error) {
	m, ok := models[modelName]
	if !ok {
		return CostEstimate{}, &ModelNotFoundError{Name: modelName}
	}
	inputCost := float64(inputTokens) / 1e6 * m.InputCostPer1M
	outputCost := float64(outputTokens) / 1e6 * m.OutputCostPer1M
	return CostEstimate{
		Model:        modelName,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		InputCost:    inputCost,
		OutputCost:   outputCost,
		TotalCost:    inputCost + outputCost,
	}, nil
}

// FitsContext checks if the given token count fits within the model's context window.
func FitsContext(modelName string, tokenCount int) (bool, int) {
	m, ok := models[modelName]
	if !ok {
		return false, 0
	}
	return tokenCount <= m.MaxTokens, m.MaxTokens - tokenCount
}

// ModelNotFoundError is returned when a model is not in the database.
type ModelNotFoundError struct {
	Name string
}

func (e *ModelNotFoundError) Error() string {
	return "model not found: " + e.Name
}
