package database

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Config for the database.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Postgres bool
	Temp     bool
}

// Connect attempt to connect to the database based on the config options.
func (c Config) Connect() (*gorm.DB, error) {

	switch {
	case c.Temp:
		return tempDB()
	case c.Postgres:
		return openPostgres(c.User, c.Password, c.Host, c.Port, c.Database)
	default:
		return openMysql(c.User, c.Password, c.Host, c.Port, c.Database)
	}

	return nil, errors.New("unable to connect to a database of any kind")

}
