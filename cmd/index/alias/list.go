package alias

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/alias/list"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List index aliases",
	Long:  `List index aliases`,
	Run: func(cmd *cobra.Command, _ []string) {
		app := app.New()
		handler := list.NewHandler(app)
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
	flags := ListCmd.Flags()
	flags.String("format", "", "Pretty-print result using a Go template")
}
