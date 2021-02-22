package index

import (
	"esctl/internal/index/app"
	"esctl/internal/index/move"
	"log"

	"github.com/spf13/cobra"
)

var MoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move src index to dest index",
	Long:  `Move src index to dest index`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := move.NewHandler(app)
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
	flags := MoveCmd.Flags()
	flags.StringP("src", "s", "", "The src index name")
	flags.StringP("dest", "d", "", "The dest index name")
	flags.BoolP("purge", "p", false, "Delete src index after migrate success")

	if err := cobra.MarkFlagRequired(flags, "src"); err != nil {
		log.Fatal(err)
	}

	if err := cobra.MarkFlagRequired(flags, "dest"); err != nil {
		log.Fatal(err)
	}
}
