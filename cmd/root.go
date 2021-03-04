/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"esctl/cmd/infra"
	"esctl/pkg/log"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "esctl",
	Short: "esctl controls the ElasticSearch cluster manager",
	Long:  `esctl controls the ElasticSearch cluster manager`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	// 全局 flasgs，对所有子命令有效
	pFlags := rootCmd.PersistentFlags()
	pFlags.StringP("config", "c", "", "config file (default is $HOME/.esctl/config)")
	pFlags.String("context", "", "The name of the config context to use")
	pFlags.String("cluster", "", "The name of the config cluster to use")
	pFlags.String("user", "", "The name of the config user to use")

	flags := rootCmd.Flags()
	flags.BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(onInitialize)
}

func onInitialize() {
	pFlags := rootCmd.PersistentFlags()

	cfgFile, err := pFlags.GetString("config")
	if err != nil {
		fmt.Println("failed to get --config: " + err.Error())
		os.Exit(1)
	}

	// 初始化配置
	cfg, err := initConfig(cfgFile, pFlags)
	if err != nil {
		fmt.Println("failed to load config: " + err.Error())
		os.Exit(1)
	}
	// 初始化日志
	if err := infra.InitLogHelper(log.HelperConfig{
		LogLevel:  cfg.Log.Level,
		LogFormat: cfg.Log.Format,
	}); err != nil {
		fmt.Println("failed to init loggerHelper: " + err.Error())
		os.Exit(1)
	}
}
