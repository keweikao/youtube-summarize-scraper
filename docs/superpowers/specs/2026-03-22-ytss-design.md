# YTSS (YouTube Summarize Scraper) — Design Spec

## Overview

A single Go binary CLI tool that reads a list of YouTube channels from a YAML config, fetches the latest N videos from each channel, downloads subtitles (or transcribes audio via whisper.cpp), generates summaries via configurable LLM backends, and saves all output locally. Already-processed videos are automatically skipped.

External tools (`yt-dlp`, `whisper.cpp`) are embedded in the binary via `go:embed` and released to a cache directory at runtime. Whisper models are downloaded on demand.

## CLI Interface

### Commands

```
ytss run                            # Read config.yaml, batch process all channels
ytss video <URL or VIDEO_ID>        # Summarize a single video
ytss channel <URL or @handle> -n 5  # Summarize latest N videos from a channel
```

### Global Flags

```
--config, -c       Config file path (default: ./config.yaml)
--output, -o       Output directory (default: ./output, overridable in config)
--llm              Override LLM backend (ollama / llamacpp / claude-api / gemini-cli)
--cookie-file      Path to cookie.txt (Netscape format)
--cookie-browser   Auto-extract cookie from browser (chrome / firefox / safari / edge)
--dry-run          List videos that would be processed without executing
--verbose, -v      Verbose logging
```

`ytss video` and `ytss channel` work standalone without config.yaml (using defaults).

## Config File

```yaml
# Output
output_dir: "./output"

# Subtitle language preferences (optional)
# If unset: detect video original language → fallback to English
preferred_languages:
  - ja
  - zh-Hant
  - en

# Default video count per channel
default_count: 5

# Whisper settings
whisper:
  model_dir: "~/.ytss/models"
  default_model: "base"              # Fallback model for unmatched languages

  language_models:                   # Language-specific model overrides
    ja: "medium"
    zh: "medium"
    en: "base"

  model_sources:                     # Download URLs (optional, defaults to HuggingFace)
    tiny: "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.bin"
    base: "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin"
    medium: "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.bin"

# Cookie settings (optional)
cookie:
  file: ""                           # Path to cookie.txt
  browser: ""                        # chrome / firefox / safari / edge

# LLM settings
llm:
  provider: "ollama"
  ollama:
    model: "llama3"
    endpoint: "http://localhost:11434"
  llamacpp:
    endpoint: "http://localhost:8080"
  claude_api:
    api_key: "${CLAUDE_API_KEY}"
    model: "claude-sonnet-4-20250514"
  gemini_cli: {}

# Summary settings
summary:
  prompt: "Please summarize the following video content in Traditional Chinese with key points..."
  max_tokens: 2000

# Channel list
channels:
  - url: "https://www.youtube.com/@channel-a"
    count: 10                        # Override default_count
  - url: "https://www.youtube.com/@channel-b"
  - url: "https://www.youtube.com/@channel-c"
```

## Output Structure

```
output/
├── @channel-a/
│   ├── 2026-03-20__dQw4w9WgXcQ__Rick_Astley_Never_Gonna_Give_You_Up/
│   │   ├── 2026-03-20__dQw4w9WgXcQ__subtitle.srt
│   │   ├── 2026-03-20__dQw4w9WgXcQ__transcription.md
│   │   └── 2026-03-20__dQw4w9WgXcQ__summary.md
│   └── 2026-03-18__abc123xyz__Another_Video_Title/
│       ├── 2026-03-18__abc123xyz__subtitle.srt
│       ├── 2026-03-18__abc123xyz__transcription.md
│       └── 2026-03-18__abc123xyz__summary.md
```

### Naming Rules

- Folder: `YYYY-MM-DD__{video_id}__{sanitized_title}`
- Files: `YYYY-MM-DD__{video_id}__{type}.{ext}`
- Date: video upload date
- Sanitized title: special characters and spaces removed, length limited
- `transcription.md`: subtitle content with SRT formatting stripped, plain text only

### Skip Detection

Check if a folder containing the video ID already exists in the output directory. If found, skip processing.

## Core Pipeline

```
┌─────────────────────────────────────────────────────────────────┐
│                         ytss run                                │
│                                                                 │
│  ┌──────────┐   ┌──────────────┐   ┌────────────────────────┐  │
│  │ Read     │──▶│ Fetch latest │──▶│ Filter: skip if video  │  │
│  │ config   │   │ N videos per │   │ ID folder exists in    │  │
│  │          │   │ channel      │   │ output directory       │  │
│  └──────────┘   └──────────────┘   └───────────┬────────────┘  │
│                                                 ▼               │
│                                    ┌────────────────────────┐  │
│                                    │ Subtitle Strategy      │  │
│                                    │                        │  │
│                                    │ 1. Preferred languages │  │
│                                    │    set in config?      │  │
│                                    │    ├─ Yes → find match │  │
│                                    │    └─ No → detect      │  │
│                                    │         original lang  │  │
│                                    │ 2. Manual subs first   │  │
│                                    │ 3. Auto subs second    │  │
│                                    │ 4. Fallback → English  │  │
│                                    │ 5. None → whisper      │  │
│                                    └───────────┬────────────┘  │
│                                                 ▼               │
│                              ┌──────────────────────────────┐  │
│                              │ Generate Output Files        │  │
│                              │                              │  │
│                              │ • subtitle.srt  (raw subs)   │  │
│                              │ • transcription.md (text)    │  │
│                              │ • summary.md (LLM summary)   │  │
│                              └──────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### Whisper Transcription Branch

1. Download audio via `yt-dlp` (smallest format: opus/m4a)
2. Select whisper model: `language_models[lang]` → `default_model` fallback
3. Auto-download model if not present (using `model_sources` URLs)
4. Transcribe with `whisper.cpp`
5. Delete audio file after transcription, keep only subtitle output

### Cookie Usage Strategy

Cookies are used **only when needed** to minimize account risk:

```
Attempt download (no cookie)
├─ Success → continue
└─ Fail with "sign in required" / "age-restricted"
   ├─ Cookie configured? → retry with cookie (once)
   │   ├─ Success → continue
   │   └─ Fail → log error, skip
   └─ No cookie? → log warning "cookie required", skip
```

Maps to `yt-dlp --cookies` / `--cookies-from-browser` flags. Usage is logged at WARN level.

## Internal Architecture

```
ytss/
├── main.go                  # CLI entry (cobra)
├── cmd/
│   ├── run.go               # ytss run
│   ├── video.go             # ytss video
│   └── channel.go           # ytss channel
├── config/
│   └── config.go            # YAML config parsing
├── fetcher/
│   └── fetcher.go           # yt-dlp: channel video list & metadata
├── subtitle/
│   └── subtitle.go          # Subtitle download, language strategy, SRT → plain text
├── transcriber/
│   └── transcriber.go       # Audio download + whisper.cpp transcription + model mgmt
├── summarizer/
│   ├── summarizer.go        # LLM interface definition
│   ├── ollama.go
│   ├── llamacpp.go
│   ├── claude.go
│   └── gemini.go
├── pipeline/
│   └── pipeline.go          # Orchestrates all modules
├── output/
│   └── output.go            # Naming rules, directory creation, skip detection
├── embedded/
│   └── embed.go             # go:embed yt-dlp & whisper.cpp, extract to cache
└── config.example.yaml
```

### Key Design Decisions

- **`summarizer` uses an interface** — all LLM backends implement `Summarize(text string) (string, error)`, easy to extend
- **`embedded/` handles binary extraction** — checks `~/.ytss/bin/` at startup, extracts from embed if missing or version mismatch
- **`pipeline/` is the single orchestration point** — all three commands call into pipeline, differing only in input source

## External Dependencies

### Embedded Binaries

| Tool | Source | Embed Strategy |
|------|--------|---------------|
| `yt-dlp` | GitHub Release (platform-specific binary) | `go:embed` |
| `whisper.cpp` | GitHub Release (`main` binary) | `go:embed` |

### Build Process

```
Makefile
├── download-deps      # Download platform-specific yt-dlp & whisper.cpp to embedded/bin/
├── build              # go build with embedded binaries
├── build-all          # Cross-compile: darwin-arm64 / darwin-amd64 / linux-amd64
└── clean
```

- Platform-specific binaries require per-platform builds
- `embedded/bin/` is in `.gitignore`; CI downloads correct versions

### Go Dependencies

| Purpose | Library |
|---------|---------|
| CLI framework | `github.com/spf13/cobra` |
| YAML parsing | `gopkg.in/yaml.v3` |
| HTTP client (LLM API) | stdlib `net/http` |
| Logging | stdlib `log/slog` |

## Error Handling & Logging

### Error Strategy

- **Single video failure does not abort batch** — log error, continue to next
- **External tool failure** — capture stderr, mark video as failed
- **LLM unavailable** — produce subtitle and transcription normally, skip summary only, log warning
- **Network issues** — on model download failure, suggest manual download path

### Log Format

Structured logging via `slog`, default info level, `-v` for debug:

```
INFO  processing channel @channel-a (5 videos)
INFO  [1/5] dQw4w9WgXcQ - skipped (already exists)
INFO  [2/5] abc123xyz - downloading subtitle (ja, manual)
INFO  [2/5] abc123xyz - generating summary (ollama/llama3)
INFO  [2/5] abc123xyz - done
WARN  [3/5] def456uvw - no subtitle available, transcribing with whisper (medium)
WARN  [3/5] def456uvw - cookie required for download, retrying with cookie
ERROR [4/5] ghi789rst - failed: yt-dlp exit code 1: video unavailable
INFO  [5/5] jkl012mno - done
INFO  completed: 3 success, 1 skipped, 1 failed
```

### Completion Summary

Batch runs print statistics at the end: success / skipped / failed counts, with video IDs and reasons for failures.
