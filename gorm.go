package growl

import (
	"errors"

	"github.com/jinzhu/gorm"
)

var connDb *gorm.DB

func Conn() (db *gorm.DB, err error) {

	if connDb == nil {
		connDb, err = dbConnect()
		return connDb, err
	} else {
		err = connDb.DB().Ping()
	}

	if err != nil {
		connDb.Close()
		connDb, err = dbConnect()
		return connDb, err
	}

	return connDb, nil
}

func dbConnect() (db *gorm.DB, err error) {
	config := YamlConfig.Growl

	gorm.DefaultTableNameHandler = dbSetPrefix

	url := config.Database.Url + config.Database.Name
	db, err = gorm.Open(config.Database.Driver, url)
	if err != nil {
		return db, errors.New("error while connecting to db : " + err.Error())
	}
	db.SingularTable(config.Database.SingularTable)
	db.LogMode(config.Misc.Log)
	return db, nil
}

func dbSetPrefix(db *gorm.DB, defaultTableName string) string {
	config := YamlConfig.Growl
	return config.Database.Prefix + defaultTableName
}
