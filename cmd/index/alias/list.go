package alias

import (
	"esctl/internal/index/alias/list"
	"esctl/internal/index/app"
	"log"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List aliases of indices",
	Long:  `List aliases of indices`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := list.NewHandler(app)
		handlerFlags, err := handler.ParseCmdFlags(cmd.Flags())
		if err != nil {
			log.Fatal(err)
		}
		if err := handler.Run(handlerFlags); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	flags := ListCmd.Flags()
	flags.String("format", "", "Pretty-print indices using a Go template")
}
