package cmd

import (
	"fmt"

	"broadcastle.co/code/icii/pkg/web"
	"github.com/spf13/cobra"
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

func init() {
	RootCmd.AddCommand(webCmd)

	webCmd.Flags().IntVarP(&port, "port", "p", 8080, "port for the server to run")
}
