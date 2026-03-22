package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var videoCmd = &cobra.Command{
	Use:   "video [URL or VIDEO_ID]",
	Short: "Summarize a single video",
	Long:  "Download subtitles, transcribe if needed, and generate a summary for a single video.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("ytss video %s: not yet implemented\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(videoCmd)
}
