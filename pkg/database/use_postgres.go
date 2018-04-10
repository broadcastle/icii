package database

import (
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // This lets us use postgres
)

func openPostgres(user string, password string, host string, port int, database string) (*gorm.DB, error) {

	u := "host=" + host + " port=" + strconv.Itoa(port) + " user=" + user + " dbname=" + database + " password=" + password

	db, err := gorm.Open("postgres", u)
	if err != nil {
		return nil, err
	}

	initIciiTables(db)

	return db, nil
}
