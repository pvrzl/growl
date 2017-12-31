package growl

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/inflection"
)

func (db Db) GetTableName() string {
	config := YamlConfig.Growl.Database
	rawSplit := strings.Split(GetStructName(db.data), ".")
	name := strings.ToLower(ToSnake(rawSplit[len(rawSplit)-1]))

	if !config.SingularTable {
		return YamlConfig.Growl.Database.Prefix + inflection.Plural(name)
	}

	return YamlConfig.Growl.Database.Prefix + name
}

func (db Db) Error() error {
	return db.error
}

func (db Db) GetTx() *gorm.DB {
	return db.tx
}

func (db Db) Begin() Db {
	db.tx, db.error = Conn()
	db.tx = db.tx.Begin()
	db.txMode = true
	return db
}

func (db Db) SetTx(tx *gorm.DB) Db {
	db.tx = tx
	if tx != nil {
		db.txMode = true
	} else {
		db.txMode = false
	}

	return db
}

// func (db Db) checkTag() error {
// 	return nil
// }

func (db Db) Commit() Db {
	db.tx.Commit()
	return db
}

func (db Db) Rollback() Db {
	db.tx.Rollback()
	return db
}

func (db Db) checkTx() (Db, *gorm.DB) {
	if db.tx != nil {
		return db, db.tx
	}
	db.tx, db.error = Conn()
	db.tx = db.tx.Begin()
	return db, db.tx
}
