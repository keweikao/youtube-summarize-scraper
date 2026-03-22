package summarizer

import (
	"fmt"
	"strings"
)

// KeywordPrompt generates a Stage 2 keyword extraction prompt.
// The prompt language is selected by the language parameter.
func KeywordPrompt(summary string, language string, count int) string {
	switch language {
	case "zh-Hant":
		return fmt.Sprintf(
			"請從以下摘要中提取最多 %d 個關鍵字，每行列出一個關鍵字，不要編號，不要其他說明文字。使用繁體中文，遇到專有名詞保留原文。\n\n%s",
			count, summary,
		)
	case "ja":
		return fmt.Sprintf(
			"以下の要約から最大 %d 個のキーワードを抽出してください。1行に1つのキーワードを記載し、番号や説明は不要です。専門用語は原語を保持してください。\n\n%s",
			count, summary,
		)
	default:
		return fmt.Sprintf(
			"Extract up to %d keywords from the summary below. List one keyword per line. No numbering, no extra text.\n\n%s",
			count, summary,
		)
	}
}

// ParseKeywords splits an LLM response into individual keywords.
// It trims whitespace, removes bullet markers (-, *, bullet, numbers), and discards empty lines.
func ParseKeywords(response string) []string {
	lines := strings.Split(response, "\n")
	var keywords []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove common bullet markers
		line = strings.TrimLeft(line, "-*\u2022 ")
		// Remove numbered list prefixes like "1.", "2)", "10."
		line = trimNumberPrefix(line)
		line = strings.TrimSpace(line)

		if line != "" {
			keywords = append(keywords, line)
		}
	}

	return keywords
}

// trimNumberPrefix removes leading number prefixes like "1.", "2)", "10." from a line.
func trimNumberPrefix(s string) string {
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i > 0 && i < len(s) && (s[i] == '.' || s[i] == ')') {
		return strings.TrimSpace(s[i+1:])
	}
	return s
}
