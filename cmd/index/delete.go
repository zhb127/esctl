package index

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/index/delete"

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
		}

		app := app.New()
		handler := delete.NewHandler(app)
		if err := handler.Run(args); err != nil {
			infra.LogHelper.Fatal("failed to run handler", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			infra.LogHelper.Info("ok", nil)
		}
	},
}

func init() {
}
