/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"esctl/internal/index/app"
	"esctl/internal/index/list"
	"log"

	"github.com/spf13/cobra"
)

// indexListCmd represents the indexList command
var indexListCmd = &cobra.Command{
	Use:   "list (NAME_WITH_WILDCARDS_1 ... NAME_WITH_WILDCARDS_N)",
	Short: "Lists the specified indices",
	Long:  `Lists the specified indices`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := list.NewHandler(app)
		handlerFlags, err := handler.ParseCmdFlags(cmd.Flags())
		if err != nil {
			log.Fatal(err)
		}
		if err := handler.Handle(handlerFlags, args); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	indexCmd.AddCommand(indexListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indexListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indexListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	indexListCmd.Flags().String("format", "", "Pretty-print indices using a Go template")
	indexListCmd.Flags().BoolP("all", "a", false, "List including hidden indices")
}
