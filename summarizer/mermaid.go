package summarizer

import (
	"fmt"
	"strings"
)

// MermaidPrompt generates a Stage 3 Mermaid flowchart prompt.
// The prompt language is selected by the language parameter.
func MermaidPrompt(summary string, language string) string {
	switch language {
	case "zh-Hant":
		return fmt.Sprintf(
			"請根據以下影片摘要，用 Mermaid 流程圖呈現影片的敘事邏輯或核心概念的關係。\n\n"+
				"規則：\n"+
				"- 使用 graph TD（上到下）格式\n"+
				"- 節點文字用雙引號包裹，例如：A[\"節點文字\"]\n"+
				"- 只用簡單箭頭 -->\n"+
				"- 節點數量控制在 5-12 個\n"+
				"- 只輸出 Mermaid 語法區塊，不要其他說明文字\n\n%s",
			summary,
		)
	case "ja":
		return fmt.Sprintf(
			"以下の動画要約に基づき、Mermaid フローチャートで動画の論理構成または核心概念の関係を表現してください。\n\n"+
				"ルール：\n"+
				"- graph TD（上から下）形式を使用\n"+
				"- ノードテキストはダブルクォートで囲む。例：A[\"ノードテキスト\"]\n"+
				"- 矢印は --> のみ使用\n"+
				"- ノード数は 5-12 個に制限\n"+
				"- Mermaid コードブロックのみ出力、説明文不要\n\n%s",
			summary,
		)
	default:
		return fmt.Sprintf(
			"Based on the video summary below, create a Mermaid flowchart showing the narrative logic or relationships between core concepts.\n\n"+
				"Rules:\n"+
				"- Use graph TD (top-down) format\n"+
				"- Wrap node text in double quotes, e.g.: A[\"Node text\"]\n"+
				"- Use only simple arrows -->\n"+
				"- Keep nodes between 5-12\n"+
				"- Output only the Mermaid code block, no other text\n\n%s",
			summary,
		)
	}
}

// ValidateMermaid extracts and validates a Mermaid code block from an LLM response.
// It looks for ```mermaid ... ``` blocks, validates basic syntax requirements,
// and returns the cleaned content and a validity flag.
func ValidateMermaid(content string) (string, bool) {
	// Extract mermaid code block
	mermaidCode := extractMermaidBlock(content)
	if mermaidCode == "" {
		return "", false
	}

	// Basic validation: must start with "graph" or "flowchart"
	trimmed := strings.TrimSpace(mermaidCode)
	if !strings.HasPrefix(trimmed, "graph") && !strings.HasPrefix(trimmed, "flowchart") {
		return "", false
	}

	// Must contain at least one arrow
	if !strings.Contains(trimmed, "-->") {
		return "", false
	}

	return trimmed, true
}

// extractMermaidBlock finds and extracts content between ```mermaid and ``` markers.
func extractMermaidBlock(content string) string {
	startMarker := "```mermaid"
	endMarker := "```"

	startIdx := strings.Index(content, startMarker)
	if startIdx == -1 {
		return ""
	}

	// Move past the start marker
	codeStart := startIdx + len(startMarker)
	// Skip any whitespace/newline after the marker
	for codeStart < len(content) && (content[codeStart] == '\n' || content[codeStart] == '\r') {
		codeStart++
	}

	// Find the closing ```
	remaining := content[codeStart:]
	endIdx := strings.Index(remaining, endMarker)
	if endIdx == -1 {
		return ""
	}

	return strings.TrimSpace(remaining[:endIdx])
}
