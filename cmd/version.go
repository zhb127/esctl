package cmd

import (
	"esctl/cmd/infra"
	"esctl/internal/app"
	"esctl/internal/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version information",
	Long:  `Show the version information`,

	Run: func(_ *cobra.Command, _ []string) {
		app := app.New()
		handler := version.NewHandler(app)
		if err := handler.Run(); err != nil {
			infra.LogHelper.Fatal(err.Error(), nil)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
