package cmd

import (
	"fmt"
	"likha/config"
	"likha/runner"
	"likha/util"

	"github.com/spf13/cobra"
)

var (
	configPath string
	countStr   string
	output     string
)

var rootCmd = &cobra.Command{
	Use:   "likha",
	Short: "A powerful tool to generate large volumes of test data.",
	Long: `likha is a flexible and performant Command Line tool for generating
test data in various formats like CSV, JSON, XML, and YAML.
It supports a wide range of value generators and allows for complex data structures.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		count, err := util.ParseCount(countStr)
		if err != nil {
			return fmt.Errorf("invalid count value: %w", err)
		}

		if output != "" {
			cfg.Output.File = output
		}

		r, err := runner.NewRunner(cfg, count)
		if err != nil {
			return fmt.Errorf("failed to initialize runner: %w", err)
		}

		// Always run with the progress bar
		return r.Run()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to the configuration file.")
	rootCmd.Flags().StringVarP(&countStr, "count", "n", "100", "Number of records to generate (e.g., 10, 10k, 10m, 1b).")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (overrides config).")
	// Removed the --progress flag as it's no longer needed
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
