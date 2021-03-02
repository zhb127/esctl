package index

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/create"

	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create index",
	Long:  `Create index`,
	Run: func(cmd *cobra.Command, _ []string) {
		app := app.New()
		handler := create.NewHandler(app)
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
	flags := CreateCmd.Flags()
	flags.String("name", "", "The name of the index to be created")
	flags.String("body", "", "The body (index mapping JSON schema) of the index to be created")
	flags.StringP("file", "f", "", "The body (index mapping JSON schema file path) of the index to be created")

	if err := cobra.MarkFlagRequired(flags, "name"); err != nil {
		infra.LogHelper.Fatal(err.Error(), nil)
	}
}
