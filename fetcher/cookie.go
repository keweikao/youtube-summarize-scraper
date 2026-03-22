package fetcher

// cookieArgs returns yt-dlp cookie arguments based on the configured cookie settings.
func (f *Fetcher) cookieArgs() []string {
	if f.cookieConfig.File != "" {
		return []string{"--cookies", f.cookieConfig.File}
	}
	if f.cookieConfig.Browser != "" {
		browser := f.cookieConfig.Browser
		if f.cookieConfig.ChromeProfile != "" {
			browser += ":" + f.cookieConfig.ChromeProfile
		}
		return []string{"--cookies-from-browser", browser}
	}
	return nil
}

// needsCookie returns true when the video's availability requires authentication cookies.
func (f *Fetcher) needsCookie(availability string) bool {
	switch availability {
	case "members_only", "needs_auth", "premium_only", "subscriber_only", "private":
		return true
	default:
		return false
	}
}
