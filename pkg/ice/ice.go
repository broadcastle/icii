package ice

import (
	"broadcastle.co/code/icii/pkg/database"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var db *gorm.DB

// Start the database.
func Start() error {

	c := database.Config{
		Temp: viper.GetBool("database.temp"),
	}

	if !c.Temp {
		c.Database = viper.GetString("database.database")
		c.Host = viper.GetString("database.host")
		c.Password = viper.GetString("database.password")
		c.Port = viper.GetInt("database.port")
		c.Postgres = viper.GetBool("database.postgres")
		c.User = viper.GetString("database.user")
	}

	var err error

	db, err = c.Connect()
	if err != nil {
		return err
	}

	return nil

}

// Close the database.
func Close() error {
	return db.Close()
}
