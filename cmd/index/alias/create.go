package alias

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/alias/create"

	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create index alias",
	Long:  `Create index alias`,
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
	flags.String("index", "", "Index name")
	flags.String("alias", "", "Alias name")

	if err := cobra.MarkFlagRequired(flags, "index"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}

	if err := cobra.MarkFlagRequired(flags, "alias"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}
}
