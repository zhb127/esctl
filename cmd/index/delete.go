package index

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/delete"
	"os"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete INDEX_NAME_1 ... INDEX_NAME_N",
	Short: "Delete the specified indices",
	Long:  `Delete the specified indices`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			infra.LogHelper.Fatal("failed to run handler", map[string]interface{}{
				"error": "Need to specify at least one index name",
			})
			os.Exit(1)
		}

		app := app.New()
		handler := delete.NewHandler(app)
		if err := handler.Run(args); err != nil {
			infra.LogHelper.Fatal("failed to run handler", map[string]interface{}{
				"error": err.Error(),
			})
		}
	},
}

func init() {
}
