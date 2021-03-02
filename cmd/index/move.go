package index

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/move"

	"github.com/spf13/cobra"
)

var MoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move source index to destination index",
	Long:  `Move source index to destination index`,
	Run: func(cmd *cobra.Command, _ []string) {
		app := app.New()
		handler := move.NewHandler(app)
		handlerFlags, err := handler.ParseCmdFlags(cmd.Flags())
		if err != nil {
			infra.LogHelper.Fatal("failed to parse cmd flags", map[string]interface{}{
				"error": err.Error(),
			})
		}
		if err := handler.Run(handlerFlags); err != nil {
			infra.LogHelper.Fatal("failed to run handler", map[string]interface{}{
				"error": err.Error(),
			})
		}
	},
}

func init() {
	flags := MoveCmd.Flags()
	flags.StringP("src", "s", "", "The source index name")
	flags.StringP("dest", "d", "", "The destination index name")
	flags.BoolP("purge", "p", false, "Delete source index after successful operation")

	if err := cobra.MarkFlagRequired(flags, "src"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}

	if err := cobra.MarkFlagRequired(flags, "dest"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}
}
