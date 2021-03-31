package cmd

import (
	"esctl/cmd/doc"

	"github.com/spf13/cobra"
)

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Manage index documents",
	Long:  `Manage index documents`,
}

func init() {
	rootCmd.AddCommand(docCmd)

	docCmd.AddCommand(doc.DeleteCmd)
}
