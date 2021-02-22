package index

import (
	"esctl/internal/index/app"
	"esctl/internal/index/create"
	"log"

	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create index",
	Long:  `Create index`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := create.NewHandler(app)
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
	flags := CreateCmd.Flags()
	flags.String("name", "", "The name of the index to be created")
	flags.String("body", "", "The body (JSON mapping) of the index to be created")
	flags.StringP("file", "f", "", "The body (JSON mapping file path) of the index to be created")

	if err := cobra.MarkFlagRequired(flags, "name"); err != nil {
		log.Fatal(err)
	}
}
