package database

import (
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // This lets us use mysql
)

func openMysql(user string, password string, host string, port int, database string) (*gorm.DB, error) {

	u := user + ":" + password + "@" + host

	if port != 0 {
		u = u + ":" + strconv.Itoa(port)
	}

	u = u + database + "?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open("mysql", u)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Song{})

	return db, nil
}
