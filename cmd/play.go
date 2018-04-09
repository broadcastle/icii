package cmd

import (
	"fmt"
	"log"

	"broadcastle.co/code/icii/pkg/stream"
	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("play called")
		if single != "" {
			log.Println(settings.Config.Play(single))
			log.Printf("%s was played\n", single)
		}
	},
}

var settings struct {
	stream.Config
}

var single string

func init() {
	RootCmd.AddCommand(playCmd)

	playCmd.Flags().StringVarP(&settings.Config.Host, "server", "s", "", "server to connect to")
	playCmd.Flags().IntVarP(&settings.Config.Port, "port", "", 8080, "port to connect to")
	playCmd.Flags().StringVarP(&settings.Config.User, "user", "u", "source", "username to use")
	playCmd.Flags().StringVarP(&settings.Config.Password, "password", "p", "hackme", "password to use")
	playCmd.Flags().StringVarP(&settings.Config.Mount, "mount", "m", "source.mp3", "mountplace to use")
	playCmd.Flags().StringVarP(&settings.Config.Name, "name", "", "Sample Station", "name for the stream")
	playCmd.Flags().StringVarP(&settings.Config.URL, "url", "", "google.com", "url displayed in the stream")
	playCmd.Flags().StringVarP(&settings.Config.Genre, "genre", "", "pop", "genre of the stream")
	playCmd.Flags().StringVarP(&settings.Config.Description, "description", "", "sample description", "description of the stream")
	playCmd.Flags().IntVarP(&settings.Config.BufferSize, "buffer", "", 3, "number of seconds to have as a buffer")
	playCmd.Flags().StringVarP(&single, "single", "", "", "a single file to play")
}
