package growl

import (
	"errors"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
)

var connDb *gorm.DB
var dbOnce sync.Once

func Conn() (db *gorm.DB, err error) {

	dbOnce.Do(func() {
		connDb, err = dbConnect()
		if err != nil {
			panic(err)
		}
	})

	return connDb, nil
}

func dbConnect() (db *gorm.DB, err error) {
	config := YamlConfig.Growl

	gorm.DefaultTableNameHandler = dbSetPrefix

	url := config.Database.Url + config.Database.Name
	db, err = gorm.Open(config.Database.Driver, url)
	if err != nil {
		newErr := errors.New("error while connecting to db : " + err.Error())
		log.Panic(newErr)
		return db, newErr
	}
	db.SingularTable(config.Database.SingularTable)
	db.LogMode(config.Misc.Log)
	// db.DB().SetConnMaxLifetime(time.Minute * 5)
	// db.DB().SetMaxIdleConns(0)
	// db.DB().SetMaxOpenConns(5)
	return db, nil
}

func dbSetPrefix(db *gorm.DB, defaultTableName string) string {
	config := YamlConfig.Growl
	return config.Database.Prefix + defaultTableName
}
