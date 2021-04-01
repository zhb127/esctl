package doc

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/doc/delete"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete index documents by query",
	Short: "Delete the specified documents",
	Long:  `Delete the specified documents`,
	Run: func(cmd *cobra.Command, _ []string) {
		app := app.New()
		handler := delete.NewHandler(app)
		handlerFlags, err := handler.ParseCmdFlags(cmd.Flags())
		if err != nil {
			infra.LogHelper.Fatal("failed to parse cmd flags", map[string]interface{}{
				"err": err,
			})
		}

		if err := handler.Run(handlerFlags); err != nil {
			infra.LogHelper.Fatal("failed to run handler", map[string]interface{}{
				"err": err,
			})
			return
		}

		infra.LogHelper.Info("ok", nil)
	},
}

func init() {
	flags := DeleteCmd.Flags()
	flags.StringP("index", "i", "", "Index name")
	if err := cobra.MarkFlagRequired(flags, "index"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}

	flags.StringP("query", "q", "", "Query DSL json string")
	flags.BoolP("all", "a", false, "Query all documents")
}
