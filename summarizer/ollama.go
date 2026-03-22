package summarizer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaSummarizer calls the Ollama REST API.
type OllamaSummarizer struct {
	endpoint string
	model    string
}

type ollamaRequest struct {
	Model   string        `json:"model"`
	Prompt  string        `json:"prompt"`
	Stream  bool          `json:"stream"`
	Options ollamaOptions `json:"options"`
}

type ollamaOptions struct {
	NumPredict int `json:"num_predict"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (o *OllamaSummarizer) Summarize(text string, opts SummarizeOptions) (string, error) {
	model := o.model
	if opts.Model != "" {
		model = opts.Model
	}

	combinedPrompt := resolvePrompt(text, opts)

	reqBody := ollamaRequest{
		Model:  model,
		Prompt: combinedPrompt,
		Stream: false,
		Options: ollamaOptions{
			NumPredict: opts.MaxTokens,
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("ollama: marshal request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	url := o.endpoint + "/api/generate"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ollama: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama: send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ollama: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var result ollamaResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("ollama: parse response: %w", err)
	}

	return result.Response, nil
}
