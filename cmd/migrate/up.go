package migrate

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/migrate/up"

	"github.com/spf13/cobra"
)

var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "Exec up migration file",
	Long:  `Exec up migration file`,
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
	upCmdflags := UpCmd.Flags()
	upCmdflags.StringP("dir", "d", "", "The migrations dir")
	upCmdflags.StringP("from", "", "", "The migration name to start")
	upCmdflags.StringP("to", "", "", "The migration name to end")
	upCmdflags.BoolP("reverse", "r", false, "Run migrate down (step by step)")

	if err := cobra.MarkFlagRequired(upCmdflags, "dir"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}
}
