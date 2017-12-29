package growl

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

func (db Db) Where(qry string, params ...interface{}) Db {
	db.where = append(db.where, dbWhereParams{
		qry:    qry,
		params: params,
	})
	return db
}

func (db Db) Select(qry string) Db {
	db.selct = qry
	return db
}

func (db Db) Save() Db {
	if _, err := valid.ValidateStruct(db.data); err != nil {
		db.error = err
		return db
	}

	var tx *gorm.DB
	var err error

	if db.txMode {
		tx = db.tx
	} else {
		tx, err = Conn()
		if err != nil {
			db.error = err
			return db
		}
	}

	tx.Close()

	return db
}
