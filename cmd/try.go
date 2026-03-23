package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(tryCmd)
}

var tryCmd = &cobra.Command{
	Use:   "try",
	Short: "Try a command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
