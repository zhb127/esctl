package cmd

import (
	"esctl/cmd/migrate"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage migrations",
	Long:  `Manage migrations`,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.AddCommand(migrate.UpCmd)
	migrateCmd.AddCommand(migrate.CreateCmd)
}
