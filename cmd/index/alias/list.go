package alias

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/alias/list"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List aliases of indices",
	Long:  `List aliases of indices`,
	Run: func(cmd *cobra.Command, _ []string) {
		app := app.New()
		handler := list.NewHandler(app)
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
	flags := ListCmd.Flags()
	flags.String("format", "", "Pretty-print indices using a Go template")
}
