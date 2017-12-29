package growl

import (
	"reflect"

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
		tx = tx.Begin()
		if err != nil {
			db.error = err
			return db
		}
	}

	if err := db.checkTag(); err != nil {
		db.error = err
		if !db.txMode {
			tx.Rollback()
		}
		return db
	}

	id := reflect.ValueOf(db.data).Elem().FieldByName("Id")
	if id.IsValid() {
		id.Set(reflect.Zero(id.Type()))
	}

	if err := tx.Save(db.data).Error; err != nil {
		if !db.txMode {
			tx.Rollback()
		}
		db.error = err
		return db
	}

	if !db.txMode {
		tx.Commit()
	}

	return db
}
