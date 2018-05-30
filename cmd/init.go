package cmd

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate a configuration file for icii",
	Long: `Generate a configuration file for icii.
	
You must use --fangless with this command, otherwise you will get an error.`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := db.create(); err != nil {
			log.Panic(err)
		}

		fmt.Printf("configuration file was created at %s/config.toml\n", db.out)
	},
}

type configCreate struct {
	host      string
	user      string
	pass      string
	port      int
	post      bool
	temp      bool
	out       string
	files     string
	filesTemp string
}

var db configCreate

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVar(&db.post, "postgres", false, "use postgress")
	initCmd.Flags().IntVar(&db.port, "port", 13306, "database port")
	initCmd.Flags().StringVar(&db.host, "host", "localhost", "host address of the database")
	initCmd.Flags().StringVarP(&db.out, "output", "o", "~/.icii", "output directory of config file")
	initCmd.Flags().StringVarP(&db.pass, "password", "p", "", "database user password")
	initCmd.Flags().StringVarP(&db.user, "user", "u", "", "database user name")
	initCmd.Flags().BoolVar(&db.temp, "temp", false, "create a temporary database")
	initCmd.Flags().StringVar(&db.files, "files", "/tmp", "location where audio files are stored")
	initCmd.Flags().StringVar(&db.filesTemp, "processed", "/tmp", "the location where files are processed")

}

func (d configCreate) create() error {

	// Get the home directory.
	if d.out == "~/.icii" {

		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		d.out = path.Join(home, ".icii")

	}

	// Create the output folder if it does not exist.
	if _, err := os.Stat(d.out); os.IsNotExist(err) {
		if err := os.MkdirAll(d.out, 0644); err != nil {
			return err
		}
	}

	//// Create a secret for JWT.
	s := sha256.New()
	s.Write([]byte(time.Now().String()))
	secret := fmt.Sprintf("%x", s.Sum(nil))

	//// Write a config file.
	content := "[icii]" +
		"\njwt = \"" + secret + "\"" +
		"\n\n[database]"

	dc := "\npostgres = " + strconv.FormatBool(d.post) +
		"\nhost = " + db.host +
		"\nuser = " + db.user +
		"\npassword = " + db.pass +
		"\nport = " + strconv.Itoa(db.port) +
		"\nname = icii"

	lc := "\n\n[files]" +
		"\nlocation = \"" + db.files + "\""

	if d.temp {
		dc = "\ntemp = true"
	}

	//// Finally, write the configuration file.
	content = content + dc + lc

	file := []byte(content)

	output := path.Join(d.out, "config.toml")

	return ioutil.WriteFile(output, file, 0644)
}
