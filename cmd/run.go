package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Batch process all channels from config",
	Long:  "Read config.yaml and process the latest N videos from each configured channel.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("ytss run: not yet implemented")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
