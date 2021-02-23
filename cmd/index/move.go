package index

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/move"

	"github.com/spf13/cobra"
)

var MoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move src index to dest index",
	Long:  `Move src index to dest index`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := move.NewHandler(app)
		handlerFlags, err := handler.ParseCmdFlags(cmd.Flags())
		if err != nil {
			infra.LogHelper.Fatal(err.Error(), nil)
		}
		if err := handler.Run(handlerFlags); err != nil {
			infra.LogHelper.Fatal(err.Error(), nil)
		}
	},
}

func init() {
	flags := MoveCmd.Flags()
	flags.StringP("src", "s", "", "The src index name")
	flags.StringP("dest", "d", "", "The dest index name")
	flags.BoolP("purge", "p", false, "Delete src index after migrate success")

	if err := cobra.MarkFlagRequired(flags, "src"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}

	if err := cobra.MarkFlagRequired(flags, "dest"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}
}
