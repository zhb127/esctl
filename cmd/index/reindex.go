package index

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/reindex"

	"github.com/spf13/cobra"
)

var ReindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Reindex source index to dest index",
	Long:  `Reindex source index to dest index`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := reindex.NewHandler(app)
		handlerFlags, err := handler.ParseCmdFlags(cmd.Flags())
		if err != nil {
			infra.LogHelper.Fatal(err.Error(), nil)
		}
		if err := handler.Handle(handlerFlags); err != nil {
			infra.LogHelper.Fatal(err.Error(), nil)
		}
	},
}

func init() {
	flags := ReindexCmd.Flags()
	flags.StringP("src", "s", "", "The source index name")
	flags.StringP("dest", "d", "", "The dest index name")

	if err := cobra.MarkFlagRequired(flags, "src"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}

	if err := cobra.MarkFlagRequired(flags, "dest"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}
}
