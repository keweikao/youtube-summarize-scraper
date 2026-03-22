package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnrichTagsForObsidian merges originalTags, keywords, a sanitized channel name,
// and autoTags into a single deduplicated list.
func EnrichTagsForObsidian(originalTags []string, keywords []string, channelName string, autoTags []string) []string {
	seen := make(map[string]struct{})
	var result []string

	add := func(tag string) {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return
		}
		lower := strings.ToLower(tag)
		if _, ok := seen[lower]; ok {
			return
		}
		seen[lower] = struct{}{}
		result = append(result, tag)
	}

	for _, t := range originalTags {
		add(t)
	}
	for _, k := range keywords {
		add(k)
	}

	// Sanitize channel name: remove @, replace spaces with hyphens, lowercase.
	sanitized := strings.TrimPrefix(channelName, "@")
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	sanitized = strings.ToLower(sanitized)
	if sanitized != "" {
		add(sanitized)
	}

	for _, t := range autoTags {
		add(t)
	}

	return result
}

// InsertWikilink inserts a wikilink reference to the transcription file
// after the frontmatter closing "---" in the summary content.
func InsertWikilink(summaryContent string, transcriptionFileName string) string {
	wikiline := fmt.Sprintf("> Full transcription: [[%s]]", transcriptionFileName)

	// Find the second "---" (closing frontmatter delimiter).
	firstIdx := strings.Index(summaryContent, "---")
	if firstIdx < 0 {
		return wikiline + "\n\n" + summaryContent
	}
	secondIdx := strings.Index(summaryContent[firstIdx+3:], "---")
	if secondIdx < 0 {
		return wikiline + "\n\n" + summaryContent
	}

	// Position right after the closing "---\n".
	insertPos := firstIdx + 3 + secondIdx + 3
	// Skip the newline after "---" if present.
	if insertPos < len(summaryContent) && summaryContent[insertPos] == '\n' {
		insertPos++
	}

	return summaryContent[:insertPos] + "\n" + wikiline + "\n" + summaryContent[insertPos:]
}

// GenerateChannelMOC creates or overwrites an _index.md file in the channel
// directory with a Dataview query listing all videos.
func GenerateChannelMOC(channelHandle string, outputDir string) error {
	channelDir := filepath.Join(outputDir, "@"+channelHandle)
	if err := EnsureDir(channelDir); err != nil {
		return fmt.Errorf("creating channel directory: %w", err)
	}

	// Use the channel directory name as the relative FROM path.
	fromPath := "@" + channelHandle

	content := fmt.Sprintf(`# @%s

`+"```dataview"+`
TABLE upload_date, duration, subtitle_type
FROM "%s"
WHERE video_id != null
SORT upload_date DESC
`+"```"+`
`, channelHandle, fromPath)

	mocPath := filepath.Join(channelDir, "_index.md")
	if err := os.WriteFile(mocPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing MOC file: %w", err)
	}

	return nil
}
