package lang

import "testing"

func TestNormalizeToISO639_1(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Special mappings (ISO 639-3 → ISO 639-1)
		{"cmn to zh", "cmn", "zh"},
		{"yue to zh", "yue", "zh"},
		{"wuu to zh", "wuu", "zh"},
		{"jpn to ja", "jpn", "ja"},
		{"kor to ko", "kor", "ko"},
		{"eng to en", "eng", "en"},
		{"fra to fr", "fra", "fr"},
		{"fre to fr", "fre", "fr"},
		{"deu to de", "deu", "de"},
		{"ger to de", "ger", "de"},
		{"spa to es", "spa", "es"},
		{"por to pt", "por", "pt"},
		{"rus to ru", "rus", "ru"},

		// Already ISO 639-1
		{"ja unchanged", "ja", "ja"},
		{"en unchanged", "en", "en"},
		{"zh unchanged", "zh", "zh"},
		{"ko unchanged", "ko", "ko"},
		{"fr unchanged", "fr", "fr"},

		// BCP-47 tags
		{"zh-Hant", "zh-Hant", "zh"},
		{"zh-Hans", "zh-Hans", "zh"},
		{"ja-JP", "ja-JP", "ja"},
		{"en-US", "en-US", "en"},
		{"en-GB", "en-GB", "en"},
		{"pt-BR", "pt-BR", "pt"},
		{"fr-CA", "fr-CA", "fr"},

		// yt-dlp special suffixes
		{"en-orig", "en-orig", "en"},
		{"en-uYU-mmqFLq8", "en-uYU-mmqFLq8", "en"},
		{"ja-orig", "ja-orig", "ja"},

		// Case insensitivity
		{"uppercase JPN", "JPN", "ja"},
		{"mixed case Eng", "Eng", "en"},
		{"uppercase ZH-HANT", "ZH-HANT", "zh"},
		{"uppercase EN-US", "EN-US", "en"},

		// Edge cases
		{"empty string", "", ""},
		{"whitespace only", "  ", ""},
		{"with leading space", " en", "en"},
		{"with trailing space", "en ", "en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeToISO639_1(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizeToISO639_1(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
