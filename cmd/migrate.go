package cmd

import (
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage migrations",
	Long:  `Manage migrations`,

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(MigrateCmd)
}
