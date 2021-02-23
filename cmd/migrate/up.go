package migrate

import (
	"esctl/internal/app"
	"esctl/internal/migrate/up"
	"log"

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
			log.Fatal(err)
		}
		if err := handler.Run(handlerFlags); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	flags := UpCmd.Flags()
	flags.StringP("dir", "d", "", "The migrations dir")
	flags.StringP("from", "", "", "File name(without ext) to start migration")
	flags.StringP("to", "", "", "File name(without ext) to end migration")

	if err := cobra.MarkFlagRequired(flags, "dir"); err != nil {
		log.Fatal(err)
	}
}
