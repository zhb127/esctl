package cmd

import (
	"esctl/internal/app"
	"esctl/internal/version"
	"log"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version information",
	Long:  `Show the version information`,

	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := version.NewHandler(app)
		if err := handler.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
