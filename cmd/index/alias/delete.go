package alias

import (
	"esctl/internal/app"
	"esctl/internal/index/alias/delete"
	"log"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete index alias",
	Long:  `Delete index alias`,
	Run: func(cmd *cobra.Command, _ []string) {
		app := app.New()
		handler := delete.NewHandler(app)
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
	flags := DeleteCmd.Flags()
	flags.String("index", "", "Index name")
	flags.String("alias", "", "Alias name")

	if err := cobra.MarkFlagRequired(flags, "index"); err != nil {
		log.Fatal(err)
	}

	if err := cobra.MarkFlagRequired(flags, "alias"); err != nil {
		log.Fatal(err)
	}
}
