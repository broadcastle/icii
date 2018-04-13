package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is the command that starts everything.
var RootCmd = &cobra.Command{
	Use:   "icii",
	Short: "broadcast to icecast using icii",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var initialize bool
var config string

func init() {

	RootCmd.PersistentFlags().BoolVar(&initialize, "fangless", false, "use this flag in conjunction with the init command")
	RootCmd.PersistentFlags().StringVar(&config, "config", "~/.icii/config.toml", "configuration file location")

	cobra.OnInitialize(initConfig)
}

// Execute is used by icii.go to start the program.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {

	if initialize {
		return
	}

	if config == "~/.icii/config.toml" {

		home, err := homedir.Dir()
		if err != nil {
			log.Panic(err)
		}

		config = path.Join(home, ".icii", "config.toml")

	}

	viper.SetConfigFile(config)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("unable to read config file: %s\ndid you run 'icii init --fangless'?\n", err.Error())
		os.Exit(1)
	}

}
