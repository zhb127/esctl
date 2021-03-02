package index

import (
	"esctl/cmd/index/alias"

	"github.com/spf13/cobra"
)

var AliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage index aliases",
	Long:  `Manage index aliases`,
}

func init() {
	AliasCmd.AddCommand(alias.ListCmd)
	AliasCmd.AddCommand(alias.CreateCmd)
	AliasCmd.AddCommand(alias.DeleteCmd)
}
