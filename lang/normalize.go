package lang

import "strings"

// specialMappings maps ISO 639-3 and other variant codes to ISO 639-1.
var specialMappings = map[string]string{
	"cmn": "zh",
	"yue": "zh",
	"wuu": "zh",
	"jpn": "ja",
	"kor": "ko",
	"eng": "en",
	"fra": "fr",
	"fre": "fr",
	"deu": "de",
	"ger": "de",
	"spa": "es",
	"por": "pt",
	"rus": "ru",
}

// NormalizeToISO639_1 converts a language code to its ISO 639-1 two-letter form.
//
// It handles ISO 639-3 codes (e.g. "jpn" → "ja"), BCP-47 tags (e.g. "zh-Hant" → "zh"),
// and yt-dlp special suffixes (e.g. "en-orig" → "en").
// If the input is empty, an empty string is returned.
func NormalizeToISO639_1(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))
	if code == "" {
		return ""
	}

	// Check special mappings first (exact match on full lowered code).
	if mapped, ok := specialMappings[code]; ok {
		return mapped
	}

	// If already two characters, return as-is.
	if len(code) == 2 {
		return code
	}

	// For longer codes, try the prefix before any separator.
	prefix := code
	if idx := strings.IndexAny(code, "-_"); idx > 0 {
		prefix = code[:idx]
	}

	// Check special mappings on the prefix.
	if mapped, ok := specialMappings[prefix]; ok {
		return mapped
	}

	// If prefix is already two characters, use it.
	if len(prefix) == 2 {
		return prefix
	}

	// Fallback: take first two characters.
	if len(code) >= 2 {
		return code[:2]
	}

	return code
}
