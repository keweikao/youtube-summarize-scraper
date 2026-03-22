package fetcher

// VideoMeta holds metadata for a single YouTube video, as returned by yt-dlp --dump-json.
type VideoMeta struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Channel        string   `json:"channel"`         // handle like @channel (uploader_id)
	ChannelName    string   `json:"channel_name"`     // display name (uploader / channel)
	UploadDate     string   `json:"upload_date"`      // YYYYMMDD
	Duration       int      `json:"duration"`         // seconds
	DurationString string   `json:"duration_string"`
	Language       string   `json:"language"`
	Tags           []string `json:"tags"`
	Categories     []string `json:"categories"`
	Availability   string   `json:"availability"`
	LiveStatus     string   `json:"live_status"`
	URL            string   `json:"webpage_url"`
}
