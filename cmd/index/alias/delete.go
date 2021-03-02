package alias

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/alias/delete"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete index alias",
	Long:  `Delete index alias`,
	Run: func(cmd *cobra.Command, _ []string) {
		app := app.New()
		handler := delete.NewHandler(app)
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
	flags := DeleteCmd.Flags()
	flags.String("index", "", "Index name")
	flags.String("alias", "", "Alias name")

	if err := cobra.MarkFlagRequired(flags, "index"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}

	if err := cobra.MarkFlagRequired(flags, "alias"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}
}
