package index

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/list"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list (WILDCARDS_1 ... WILDCARDS_N)",
	Short: "Lists the specified indices",
	Long:  `Lists the specified indices`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := list.NewHandler(app)
		handlerFlags, err := handler.ParseCmdFlags(cmd.Flags())
		if err != nil {
			infra.LogHelper.Fatal(err.Error(), nil)
		}
		if err := handler.Run(handlerFlags, args); err != nil {
			infra.LogHelper.Fatal(err.Error(), nil)
		}
	},
}

func init() {
	flags := ListCmd.Flags()
	flags.String("format", "", "Pretty-print indices using a Go template")
	flags.BoolP("all", "a", false, "List including hidden indices")
}
