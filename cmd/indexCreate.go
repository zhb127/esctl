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
	"esctl/internal/index/create"
	"log"

	"github.com/spf13/cobra"
)

// indexCreateCmd represents the indexCreate command
var indexCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create index",
	Long:  `Create index`,
	Run: func(cmd *cobra.Command, args []string) {
		app := app.New()
		handler := create.NewHandler(app)
		handlerFlags, err := handler.ParseCmdFlags(cmd.Flags())
		if err != nil {
			log.Fatal(err)
		}
		if err := handler.Handle(handlerFlags); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	indexCmd.AddCommand(indexCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indexCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indexCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	flags := indexCreateCmd.Flags()
	flags.String("name", "", "The name of the index to be created")
	flags.String("body", "", "The body (JSON) of the index to be created")
	flags.String("file", "f", "The body (JSON file path) of the index to be created")

	if err := cobra.MarkFlagRequired(flags, "name"); err != nil {
		log.Fatal(err)
	}
}
