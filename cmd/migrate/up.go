package migrate

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/migrate/up"

	"github.com/spf13/cobra"
)

var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run migrate up",
	Long:  `Run migrate up`,
	Run: func(cmd *cobra.Command, _ []string) {
		app := app.New()
		handler := up.NewHandler(app)
		handlerFlags, err := handler.ParseCmdFlags(cmd.Flags())
		if err != nil {
			infra.LogHelper.Fatal(err.Error(), nil)
		}
		if err := handler.Run(handlerFlags); err != nil {
			infra.LogHelper.Fatal(err.Error(), nil)
		}
		infra.LogHelper.Info("success", nil)
	},
}

func init() {
	flags := UpCmd.Flags()
	flags.StringP("dir", "d", "", "The migrations dir")
	flags.StringP("from", "", "", "File name(without ext) to start migration")
	flags.StringP("to", "", "", "File name(without ext) to end migration")

	if err := cobra.MarkFlagRequired(flags, "dir"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}
}
