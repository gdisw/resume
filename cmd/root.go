package cmd

import (
	"os"

	"github.com/gdisw/resume/pkg/env"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "core",
	Short: "Commands powering Resume",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := env.Load(); err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
