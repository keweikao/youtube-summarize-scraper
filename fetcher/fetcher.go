package fetcher

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/kouko/youtube-summarize-scraper/config"
)

const metadataTimeout = 60 * time.Second

// Fetcher wraps yt-dlp for retrieving video metadata and channel listings.
type Fetcher struct {
	ytdlpPath    string
	cookieConfig config.CookieConfig
}

// NewFetcher creates a Fetcher with the given yt-dlp binary path and cookie configuration.
func NewFetcher(ytdlpPath string, cookie config.CookieConfig) *Fetcher {
	return &Fetcher{
		ytdlpPath:    ytdlpPath,
		cookieConfig: cookie,
	}
}

// FetchVideoMeta retrieves full metadata for a single video URL.
// If the video requires authentication (based on availability), cookies are used automatically.
func (f *Fetcher) FetchVideoMeta(videoURL string) (*VideoMeta, error) {
	args := []string{"--dump-json", "--no-download", videoURL}
	out, err := f.runYtDlp(args, false)
	if err != nil {
		return nil, fmt.Errorf("fetching video metadata: %w", err)
	}

	var meta VideoMeta
	if err := json.Unmarshal(out, &meta); err != nil {
		return nil, fmt.Errorf("parsing video metadata JSON: %w", err)
	}

	// Retry with cookies if the video requires authentication.
	if f.needsCookie(meta.Availability) {
		out, err = f.runYtDlp(args, true)
		if err != nil {
			return nil, fmt.Errorf("fetching video metadata with cookies: %w", err)
		}
		if err := json.Unmarshal(out, &meta); err != nil {
			return nil, fmt.Errorf("parsing video metadata JSON (cookie retry): %w", err)
		}
	}

	return &meta, nil
}

// FetchChannelVideos lists videos from a channel URL, returning up to limit items.
// It uses --flat-playlist to avoid downloading full metadata for each video.
func (f *Fetcher) FetchChannelVideos(channelURL string, limit int) ([]VideoMeta, error) {
	channelVideosURL := channelURL + "/videos"
	args := []string{
		"--flat-playlist",
		"--dump-json",
		channelVideosURL,
	}

	out, err := f.runYtDlp(args, false)
	if err != nil {
		return nil, fmt.Errorf("fetching channel videos: %w", err)
	}

	var videos []VideoMeta
	scanner := bufio.NewScanner(bytes.NewReader(out))
	// Increase scanner buffer for large JSON lines.
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)
	for scanner.Scan() {
		if limit > 0 && len(videos) >= limit {
			break
		}
		var meta VideoMeta
		if err := json.Unmarshal(scanner.Bytes(), &meta); err != nil {
			continue // skip malformed lines
		}
		videos = append(videos, meta)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading channel video listing: %w", err)
	}

	return videos, nil
}

// runYtDlp executes yt-dlp with the given arguments and an optional cookie flag.
// A context timeout of 60 seconds is applied.
func (f *Fetcher) runYtDlp(args []string, useCookie bool) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), metadataTimeout)
	defer cancel()

	fullArgs := make([]string, 0, len(args)+4)
	if useCookie {
		fullArgs = append(fullArgs, f.cookieArgs()...)
	}
	fullArgs = append(fullArgs, args...)

	cmd := exec.CommandContext(ctx, f.ytdlpPath, fullArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("yt-dlp %v: %w\nstderr: %s", args, err, stderr.String())
	}

	return stdout.Bytes(), nil
}
