package index

import (
	"esctl/cmd/index/alias"

	"github.com/spf13/cobra"
)

var AliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage ES index alias",
	Long:  `Manage ES index alias`,
}

func init() {
	AliasCmd.AddCommand(alias.ListCmd)
	AliasCmd.AddCommand(alias.CreateCmd)
	AliasCmd.AddCommand(alias.DeleteCmd)
}
