package summarizer

import (
	"fmt"

	"github.com/kouko/youtube-summarize-scraper/config"
)

// SummarizeOptions holds options for a summarization request.
type SummarizeOptions struct {
	Prompt    string
	MaxTokens int
	Model     string
}

// Summarizer is the interface that all LLM backends must implement.
type Summarizer interface {
	Summarize(text string, opts SummarizeOptions) (string, error)
}

// NewSummarizer creates a Summarizer backend based on the provider in cfg.
func NewSummarizer(cfg config.LLMConfig) (Summarizer, error) {
	switch cfg.Provider {
	case "ollama":
		return &OllamaSummarizer{
			endpoint: cfg.Ollama.Endpoint,
			model:    cfg.Ollama.Model,
		}, nil
	case "llamacpp":
		return &LlamaCppSummarizer{
			endpoint: cfg.LlamaCpp.Endpoint,
		}, nil
	case "claude-api":
		return &ClaudeSummarizer{
			apiKey: cfg.ClaudeAPI.APIKey,
			model:  cfg.ClaudeAPI.Model,
		}, nil
	case "gemini-cli":
		return &GeminiCLISummarizer{
			model:      cfg.GeminiCLI.Model,
			binaryPath: cfg.GeminiCLI.Path,
		}, nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %q", cfg.Provider)
	}
}

// resolvePrompt returns opts.Prompt if non-empty, otherwise falls back to text.
func resolvePrompt(text string, opts SummarizeOptions) string {
	if opts.Prompt != "" {
		return opts.Prompt
	}
	return text
}
