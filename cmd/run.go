package cmd

import (
	"fmt"

	"github.com/kouko/youtube-summarize-scraper/pipeline"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Batch process all channels from config",
	Long:  "Read config.yaml and process the latest N videos from each configured channel.",
	RunE: func(cmd *cobra.Command, args []string) error {
		setupLogging(verbose)

		cfg := loadConfig(cfgFile)
		applyOverrides(cfg)

		if len(cfg.Channels) == 0 {
			return fmt.Errorf("no channels configured in %s", cfgFile)
		}

		p, err := pipeline.NewPipeline(cfg, forceFlag, dryRun)
		if err != nil {
			return fmt.Errorf("initializing pipeline: %w", err)
		}

		stats, err := p.ProcessBatch()
		if err != nil {
			return fmt.Errorf("batch processing: %w", err)
		}

		printStats(stats)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
