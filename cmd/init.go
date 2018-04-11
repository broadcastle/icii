package cmd

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate a configuration file for icii",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db.create(); err != nil {
			log.Panic(err)
		}

		fmt.Printf("config exists at %s/config.toml\n", db.out)
	},
}

type configCreate struct {
	host string
	user string
	pass string
	port int
	post bool
	temp bool
	out  string
}

var db configCreate

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVar(&db.post, "postgres", false, "use postgress")
	initCmd.Flags().IntVar(&db.port, "port", 13306, "database port")
	initCmd.Flags().StringVar(&db.host, "host", "localhost", "host address of the database")
	initCmd.Flags().StringVarP(&db.out, "output", "o", "/tmp", "output location of config file")
	initCmd.Flags().StringVarP(&db.pass, "password", "p", "", "database user password")
	initCmd.Flags().StringVarP(&db.user, "user", "u", "", "database user name")
	initCmd.Flags().BoolVar(&db.temp, "temp", false, "create a temporary database")

}

func (d configCreate) create() error {

	s := sha256.New()
	s.Write([]byte(time.Now().String()))

	secret := fmt.Sprintf("%x", s.Sum(nil))

	content := "[icii]" +
		"\njwt = \"" + secret + "\"" +
		"\n\n[database]"

	dc := "\npostgres = " + strconv.FormatBool(d.post) +
		"\nhost = " + db.host +
		"\nuser = " + db.user +
		"\npassword = " + db.pass +
		"\nport = " + strconv.Itoa(db.port) +
		"\nname = icii"

	if d.temp {
		dc = "\ntemp = true"
	}

	content = content + dc

	file := []byte(content)

	output := path.Join(d.out, "config.toml")

	return ioutil.WriteFile(output, file, 0644)
}
