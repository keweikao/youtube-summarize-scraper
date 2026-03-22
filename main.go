package main

import (
	"os"

	"github.com/kouko/youtube-summarize-scraper/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
