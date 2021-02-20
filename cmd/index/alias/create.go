package alias

import (
	"esctl/internal/index/alias/create"
	"esctl/internal/index/app"
	"log"

	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create index alias",
	Long:  `Create index alias`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := create.NewHandler(app)
		handlerFlags, err := handler.ParseCmdFlags(cmd.Flags())
		if err != nil {
			log.Fatal(err)
		}
		if err := handler.Handle(handlerFlags); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	flags := CreateCmd.Flags()
	flags.String("index", "", "Index name")
	flags.String("alias", "", "Alias name")

	if err := cobra.MarkFlagRequired(flags, "index"); err != nil {
		log.Fatal(err)
	}

	if err := cobra.MarkFlagRequired(flags, "alias"); err != nil {
		log.Fatal(err)
	}
}
