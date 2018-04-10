package database

import (
	"io/ioutil"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // This lets us use sqlite3
)

func tempDB() (*gorm.DB, error) {

	tmp, err := ioutil.TempFile("", "database.db")
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("sqlite3", tmp.Name())
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Song{})

	return db, nil

}
