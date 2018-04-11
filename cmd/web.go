package cmd

import (
	"fmt"
	"os"
	"path"

	"broadcastle.co/code/icii/pkg/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("web called")

		web.Start(port)

	},
}

var port int
var config string
var temp bool

func init() {

	RootCmd.AddCommand(webCmd)

	webCmd.PersistentFlags().StringVar(&config, "config", "~/.icii/config.toml", "location of the config file")
	webCmd.PersistentFlags().BoolVar(&temp, "temp", false, "run a temporary version of icii")

	webCmd.Flags().IntVarP(&port, "port", "p", 8080, "port for the server to run")
	viper.SetDefault("icii.port", port)

	cobra.OnInitialize(initConfig)
}

func initConfig() {

	if temp {

		d := configCreate{temp: true, out: "/tmp"}
		if err := d.create(); err != nil {
			fmt.Printf("unable to create temporary file: %s", err.Error())
			os.Exit(1)
		}
		config = path.Join(d.out, "config.toml")

	}

	viper.SetConfigFile(config)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("unable to read config file: %s\ndid you run 'icii init'?\n", err.Error())
		os.Exit(1)
	}

}
