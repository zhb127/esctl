package index

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/delete"
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
				infra.LogHelper.Fatal(err.Error(), nil)
			}
			os.Exit(1)
		}

		app := app.New()
		handler := delete.NewHandler(app)
		if err := handler.Run(args); err != nil {
			infra.LogHelper.Fatal(err.Error(), nil)
		}
	},
}

func init() {
}
