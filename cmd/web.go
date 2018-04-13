package cmd

import (
	"fmt"
	"os"

	"broadcastle.co/code/icii/pkg/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "start the server for icii",
	Run: func(cmd *cobra.Command, args []string) {

		if initialize {
			fmt.Println("do not use --fangless with this command")
			os.Exit(1)
		}

		web.Start(port)

	},
}

var port int

func init() {

	RootCmd.AddCommand(webCmd)

	webCmd.Flags().IntVarP(&port, "port", "p", 8080, "port for the server to run")

	viper.SetDefault("icii.port", port)

}
