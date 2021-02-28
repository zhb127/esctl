package migrate

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/migrate/create"

	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create migration up/down files",
	Long:  `Create migration up/down files`,
	Run: func(cmd *cobra.Command, _ []string) {
		app := app.New()
		handler := create.NewHandler(app)
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
	flags := CreateCmd.Flags()
	flags.StringP("dir", "d", "", "The migrations dir")
	flags.StringP("name", "n", "", "The migration name")

	if err := cobra.MarkFlagRequired(flags, "dir"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}

	if err := cobra.MarkFlagRequired(flags, "name"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}
}
