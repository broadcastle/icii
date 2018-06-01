package ice

import (
	"broadcastle.co/code/icii/pkg/database"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

var db *gorm.DB

///////////////
// Interface //
///////////////

// Data is the interface.
type Data interface {
	Create() error
	Delete() error
	Get() error
	Update(interface{}) error
	Echo(echo.Context) error
}

// New creates d.
func New(d Data) error {
	return d.Create()
}

// Remove deletes d.
func Remove(d Data) error {
	return d.Delete()
}

// Find gets the rest of d.
func Find(d Data) error {
	return d.Get()
}

// Update d with the data from i.
func Update(d Data, i interface{}) error {
	return d.Update(i)
}

// Echo gets d from c.
func Echo(d Data, c echo.Context) error {
	return d.Echo(c)
}

//////////////////
// Init returns //
//////////////////

// InitUser is used to create a empty user variable.
func InitUser() Data {
	return &User{}
}

// InitStation is used the create a empty station variable.
func InitStation() Data {
	return &Station{}
}

// InitTrack is used to create a empty track variable.
func InitTrack() Data {
	return &Track{}
}

// InitStream is used to create a empty stream variable.
func InitStream() Data {
	return &Stream{}
}

///////////////////////
// Database Controls //
///////////////////////

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
