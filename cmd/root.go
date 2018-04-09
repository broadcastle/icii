package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is the command that starts everything.
var RootCmd = &cobra.Command{
	Use:   "icii",
	Short: "broadcast to icecast using icii",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute is used by icii.go to start the program.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
