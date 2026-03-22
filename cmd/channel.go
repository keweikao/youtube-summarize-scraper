package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	channelCount int
)

var channelCmd = &cobra.Command{
	Use:   "channel [URL or @handle]",
	Short: "Summarize latest N videos from a channel",
	Long:  "Fetch the latest N videos from a YouTube channel and generate summaries for each.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("ytss channel %s -n %d: not yet implemented\n", args[0], channelCount)
		return nil
	},
}

func init() {
	channelCmd.Flags().IntVarP(&channelCount, "count", "n", 5, "number of latest videos to process")
	rootCmd.AddCommand(channelCmd)
}
