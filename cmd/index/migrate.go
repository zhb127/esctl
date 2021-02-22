package index

import (
	"esctl/internal/index/app"
	"esctl/internal/index/migrate"
	"log"

	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Create migration",
	Long:  `Create migration`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := migrate.NewHandler(app)
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
	flags := MigrateCmd.Flags()
	flags.StringP("src", "s", "", "The src index name")
	flags.StringP("dest", "d", "", "The dest index name")

	if err := cobra.MarkFlagRequired(flags, "src"); err != nil {
		log.Fatal(err)
	}

	if err := cobra.MarkFlagRequired(flags, "dest"); err != nil {
		log.Fatal(err)
	}
}
