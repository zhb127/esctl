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
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	defaults "github.com/mcuadros/go-defaults"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfg config
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "esctl",
	Short: "esctl controls the Elasticsearch cluster manager",
	Long:  `esctl controls the Elasticsearch cluster manager`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// 初始化完成后才执行 OnInitialize
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	pFlags := rootCmd.PersistentFlags()
	pFlags.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.esctl/config)")
	pFlags.StringP("context", "", "", "The name of the config context to use")
	pFlags.StringP("cluster", "", "", "The name of the config cluster to use")
	pFlags.StringP("user", "", "", "The name of the config user to use")
	pFlags.BoolP("verbose", "v", false, "Verbose output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.SetOutput(os.Stderr)

	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		cfgFile = home + "/.esctl/config"
	}

	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}

	defaults.SetDefaults(&cfg)

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(errors.Wrap(err, "Unmarshal config file"))
	}

	if err := injectFlagsToConfig(&cfg); err != nil {
		log.Fatal(errors.Wrap(err, "Inject flags to config"))
	}

	if err := validateConfig(&cfg); err != nil {
		log.Fatal(errors.Wrap(err, "Validate config"))
	}

	injectConfigToEnvVars(&cfg)
}
