package cmd

import (
	"esctl/cmd/index"

	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Manage indices",
	Long:  `Manage indices`,
}

func init() {
	rootCmd.AddCommand(indexCmd)

	indexCmd.AddCommand(index.ListCmd)
	indexCmd.AddCommand(index.CreateCmd)
	indexCmd.AddCommand(index.DeleteCmd)
	indexCmd.AddCommand(index.ReindexCmd)
	indexCmd.AddCommand(index.AliasCmd)

	indexCmd.AddCommand(index.MoveCmd)
}
