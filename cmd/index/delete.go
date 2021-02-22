package index

import (
	"esctl/internal/app"
	"esctl/internal/index/delete"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete NAME_1 ... NAME_N",
	Short: "Delete the specified indices",
	Long:  `Delete the specified indices`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}
			os.Exit(1)
		}

		app := app.New()
		handler := delete.NewHandler(app)
		if err := handler.Run(args); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
}
