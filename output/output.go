package output

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

// ledgerFile is the name of the JSON ledger stored in each channel/playlist directory.
const ledgerFile = ".ytss-processed.json"

// LedgerEntry records a processed video's location.
type LedgerEntry struct {
	Dir         string `json:"dir"`
	ProcessedAt string `json:"processed_at,omitempty"`
}

// Ledger maps videoID → LedgerEntry. One ledger per channel directory.
type Ledger map[string]LedgerEntry

// ledgerCache avoids re-reading the same ledger file within a single run.
var (
	ledgerCache   = map[string]Ledger{}
	ledgerCacheMu sync.Mutex
)

// ClearLedgerCache resets the in-memory ledger cache. Used in tests.
func ClearLedgerCache() {
	ledgerCacheMu.Lock()
	defer ledgerCacheMu.Unlock()
	ledgerCache = map[string]Ledger{}
}

// ReadLedger loads the ledger from a channel/playlist directory.
// Returns an empty Ledger (not nil) if the file doesn't exist.
func ReadLedger(channelDir string) Ledger {
	ledgerCacheMu.Lock()
	defer ledgerCacheMu.Unlock()

	if cached, ok := ledgerCache[channelDir]; ok {
		return cached
	}

	ledger := Ledger{}
	data, err := os.ReadFile(filepath.Join(channelDir, ledgerFile))
	if err == nil {
		_ = json.Unmarshal(data, &ledger)
	}
	ledgerCache[channelDir] = ledger
	return ledger
}

// WriteLedger persists the ledger to disk and updates the cache.
func WriteLedger(channelDir string, ledger Ledger) error {
	ledgerCacheMu.Lock()
	defer ledgerCacheMu.Unlock()

	data, err := json.MarshalIndent(ledger, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling ledger: %w", err)
	}
	if err := os.WriteFile(filepath.Join(channelDir, ledgerFile), data, 0o644); err != nil {
		return fmt.Errorf("writing ledger: %w", err)
	}
	ledgerCache[channelDir] = ledger
	return nil
}

// RecordProcessed adds a video to the channel's ledger.
func RecordProcessed(channelDir, videoID, videoSubDir, processedAt string) error {
	ledger := ReadLedger(channelDir)
	ledger[videoID] = LedgerEntry{Dir: videoSubDir, ProcessedAt: processedAt}
	return WriteLedger(channelDir, ledger)
}

// reVideoIDInDirName extracts a video ID from a directory name formatted as
// "YYYY-MM-DD__videoID__title" (legacy naming convention).
var reVideoIDInDirName = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}__([^_]+)__`)

// rebuildLedgerFromFrontmatter scans all subdirectories for summary.md files,
// reads the video_id from YAML frontmatter, and rebuilds the ledger.
// Falls back to extracting video ID from directory name (legacy format).
func rebuildLedgerFromFrontmatter(channelDir string) Ledger {
	ledger := Ledger{}
	entries, err := os.ReadDir(channelDir)
	if err != nil {
		return ledger
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		subDir := entry.Name()
		// Look for any summary.md file in the subdirectory.
		matches, _ := filepath.Glob(filepath.Join(channelDir, subDir, "*summary.md"))
		if len(matches) == 0 {
			continue
		}
		// Try 1: Read video_id from frontmatter.
		videoID := extractVideoIDFromFrontmatter(matches[0])
		// Try 2: Extract from directory name (legacy: YYYY-MM-DD__videoID__title).
		if videoID == "" {
			if m := reVideoIDInDirName.FindStringSubmatch(subDir); len(m) > 1 {
				videoID = m[1]
			}
		}
		if videoID != "" {
			ledger[videoID] = LedgerEntry{Dir: subDir}
		}
	}
	// Persist the rebuilt ledger.
	_ = WriteLedger(channelDir, ledger)
	return ledger
}

// extractVideoIDFromFrontmatter reads the video_id field from a markdown file's
// YAML frontmatter (between --- delimiters).
func extractVideoIDFromFrontmatter(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			if inFrontmatter {
				return "" // end of frontmatter, not found
			}
			inFrontmatter = true
			continue
		}
		if inFrontmatter && strings.HasPrefix(line, "video_id:") {
			val := strings.TrimPrefix(line, "video_id:")
			val = strings.TrimSpace(val)
			val = strings.Trim(val, "\"")
			return val
		}
	}
	return ""
}

// reUnsafe matches characters that are not alphanumeric, CJK, or underscores.
var reUnsafe = regexp.MustCompile(`[^\p{L}\p{N}_\s]`)

// reMultiUnderscore collapses consecutive underscores.
var reMultiUnderscore = regexp.MustCompile(`_+`)

// SanitizeTitle removes special characters (keeping alphanumeric, CJK chars,
// underscores), replaces spaces with underscores, and limits length to maxLen.
// If maxLen is 0, it defaults to 80. Leading/trailing underscores are stripped.
func SanitizeTitle(title string, maxLen int) string {
	if maxLen <= 0 {
		maxLen = 80
	}

	// Remove unsafe characters (keep letters, numbers, underscores, spaces).
	s := reUnsafe.ReplaceAllString(title, "")

	// Replace whitespace with underscores.
	s = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '_'
		}
		return r
	}, s)

	// Collapse consecutive underscores.
	s = reMultiUnderscore.ReplaceAllString(s, "_")

	// Trim leading/trailing underscores.
	s = strings.Trim(s, "_")

	// Truncate to maxLen (rune-aware to preserve CJK characters).
	runes := []rune(s)
	if len(runes) > maxLen {
		runes = runes[:maxLen]
	}
	s = string(runes)

	// Trim trailing underscores after truncation.
	s = strings.TrimRight(s, "_")

	return s
}

// formatDate converts a yt-dlp date string (YYYYMMDD) to YYYY-MM-DD.
// If the input is not exactly 8 characters, it is returned as-is.
func formatDate(date string) string {
	if len(date) == 8 {
		return date[:4] + "-" + date[4:6] + "-" + date[6:8]
	}
	return date
}

// VideoDir returns the output directory path for a single video:
//
//	outputDir/@channelHandle/YYYY-MM-DD__videoID__sanitizedTitle
func VideoDir(outputDir, channelHandle, uploadDate, videoID, title string) string {
	sanitized := SanitizeTitle(title, 0)
	formattedDate := formatDate(uploadDate)
	dirName := fmt.Sprintf("%s__%s__%s", formattedDate, videoID, sanitized)
	return filepath.Join(outputDir, "@"+channelHandle, dirName)
}

// VideoFilePrefix returns the file name prefix for video output files:
//
//	YYYY-MM-DD__videoID__
func VideoFilePrefix(uploadDate, videoID string) string {
	return fmt.Sprintf("%s__%s__", formatDate(uploadDate), videoID)
}

// IsProcessed checks whether a video is fully processed by consulting the
// channel ledger first, then falling back to frontmatter scan if the ledger
// is empty or missing.
func IsProcessed(outputDir, channelHandle, videoID string) (bool, error) {
	channelDir := filepath.Join(outputDir, "@"+channelHandle)
	ledger := ReadLedger(channelDir)

	// If ledger is empty, attempt rebuild from existing frontmatter.
	if len(ledger) == 0 {
		ledger = rebuildLedgerFromFrontmatter(channelDir)
	}

	if entry, ok := ledger[videoID]; ok {
		// Verify the directory still exists.
		dir := filepath.Join(channelDir, entry.Dir)
		if _, err := os.Stat(dir); err == nil {
			return true, nil
		}
		// Directory gone — remove stale entry.
		delete(ledger, videoID)
		_ = WriteLedger(channelDir, ledger)
	}

	return false, nil
}

// IsProcessedGlobal checks whether a video is fully processed in any
// channel/playlist directory under outputDir.
func IsProcessedGlobal(outputDir, videoID string) (bool, error) {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return false, nil
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		subDir := filepath.Join(outputDir, entry.Name())
		ledger := ReadLedger(subDir)
		if len(ledger) == 0 {
			ledger = rebuildLedgerFromFrontmatter(subDir)
		}
		if _, ok := ledger[videoID]; ok {
			return true, nil
		}
	}
	return false, nil
}

// FindVideoDir returns the existing video directory path for a given video ID
// by searching ledgers across all channel directories, then falling back to
// glob matching on legacy directory names.
func FindVideoDir(outputDir, videoID string) string {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		channelDir := filepath.Join(outputDir, entry.Name())
		ledger := ReadLedger(channelDir)
		if len(ledger) == 0 {
			ledger = rebuildLedgerFromFrontmatter(channelDir)
		}
		if le, ok := ledger[videoID]; ok {
			dir := filepath.Join(channelDir, le.Dir)
			if _, err := os.Stat(dir); err == nil {
				return dir
			}
		}
	}
	// Fallback: glob for legacy directory naming (YYYY-MM-DD__videoID__*).
	pattern := filepath.Join(outputDir, "*", fmt.Sprintf("*__%s__*", videoID))
	matches, _ := filepath.Glob(pattern)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

// HasFile checks whether a file matching the given suffix exists in videoDir.
// Suffix examples: "summary.md", "transcription.md", "subtitle.srt"
func HasFile(videoDir, suffix string) bool {
	pattern := filepath.Join(videoDir, "*__"+suffix)
	matches, _ := filepath.Glob(pattern)
	return len(matches) > 0
}

// PlaylistDir returns the output directory path for a playlist:
//
//	outputDir/_playlist__playlistID__sanitizedName
func PlaylistDir(outputDir, playlistID, playlistName string) string {
	sanitized := SanitizeTitle(playlistName, 0)
	dirName := fmt.Sprintf("_playlist__%s__%s", playlistID, sanitized)
	return filepath.Join(outputDir, dirName)
}

// EnsureDir creates the directory (and all parents) if it does not exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}
