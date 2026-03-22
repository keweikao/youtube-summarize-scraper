package summarizer

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// GeminiCLISummarizer invokes the Gemini CLI tool.
type GeminiCLISummarizer struct {
	model      string
	binaryPath string
}

func (g *GeminiCLISummarizer) Summarize(text string, opts SummarizeOptions) (string, error) {
	model := g.model
	if opts.Model != "" {
		model = opts.Model
	}

	combinedPrompt := resolvePrompt(text, opts)

	binary := g.binaryPath
	if binary == "" {
		var err error
		binary, err = exec.LookPath("gemini")
		if err != nil {
			return "", fmt.Errorf("gemini-cli: binary not found in PATH: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	args := []string{"-m", model, "-p", combinedPrompt}
	cmd := exec.CommandContext(ctx, binary, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(combinedPrompt)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("gemini-cli: execution failed: %w\nstderr: %s", err, stderr.String())
	}

	return StripThinkingTags(strings.TrimSpace(stdout.String())), nil
}
