package growl

import (
	"reflect"

	valid "github.com/asaskevich/govalidator"
)

func (db Db) Where(qry string, params ...interface{}) Db {
	db.where = append(db.where, dbWhereParams{
		qry:    qry,
		params: params,
	})
	db, tx := db.checkTx()
	db.tx = tx.Where(qry, params...)
	return db
}

func (db Db) Select(qry string, filters ...func(string) string) Db {
	for _, filter := range filters {
		qry = filter(qry)
	}
	db.selct = qry
	db, tx := db.checkTx()
	db.tx = tx.Select(qry)
	return db
}

func (db Db) Preload(qry string) Db {
	db.preload = append(db.preload, qry)
	db, tx := db.checkTx()
	db.tx = tx.Preload(qry)
	return db
}

func (db Db) Join(qry string, params ...interface{}) Db {
	db.join = append(db.join, dbWhereParams{
		qry:    qry,
		params: params,
	})
	db, tx := db.checkTx()
	db.tx = tx.Joins(qry, params...)
	return db
}

func (db Db) Association(qry string) Db {
	db.association = qry
	return db
}

func (db Db) OrderBy(qry string, filters ...func(string) string) Db {
	for _, filter := range filters {
		qry = filter(qry)
	}
	db.orderBy = qry
	db, tx := db.checkTx()
	db.tx = tx.Order(qry)
	return db
}

func (db Db) Limit(qry int, filters ...func(int) int) Db {
	for _, filter := range filters {
		qry = filter(qry)
	}
	db.limit = qry
	db, tx := db.checkTx()
	db.tx = tx.Limit(qry)
	return db
}

func (db Db) Offset(qry int, filters ...func(int) int) Db {
	for _, filter := range filters {
		qry = filter(qry)
	}
	db.offset = qry
	db, tx := db.checkTx()
	db.tx = tx.Offset(qry)
	return db
}

func (db Db) Replace(data interface{}) Db {
	db, tx := db.checkTx()
	db.error = tx.Model(db.data).Association(db.association).Replace(data).Error
	if !db.txMode {
		tx.Commit()
	}

	return db
}

func (db Db) Save() Db {
	if _, err := valid.ValidateStruct(db.data); err != nil {
		db.error = err
		return db
	}

	db, tx := db.checkTx()

	db = db.checkTag()
	if db.error != nil {
		if !db.txMode {
			tx.Rollback()
		}
		return db
	}

	id := reflect.ValueOf(db.data).Elem().FieldByName("Id")
	if id.IsValid() {
		id.Set(reflect.Zero(id.Type()))
	}

	if err := tx.Create(db.data).Error; err != nil {
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

func (db Db) ForceUpdate() Db {
	if _, err := valid.ValidateStruct(db.data); err != nil {
		db.error = err
		return db
	}

	db, tx := db.checkTx()

	db = db.checkTag()
	if db.error != nil {
		if !db.txMode {
			tx.Rollback()
		}
		return db
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

func (db Db) Update(data map[string]interface{}) Db {

	db, tx := db.checkTx()

	if err := tx.Update(data).Error; err != nil {
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

func (db Db) First() Db {
	db, tx := db.checkTx()
	if err := tx.First(db.data).Error; err != nil {
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

func (db Db) Find(data interface{}) Db {
	db, tx := db.checkTx()
	if err := tx.Find(data).Error; err != nil {
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

func (db Db) Count(data interface{}) Db {
	db, tx := db.checkTx()
	if err := tx.Count(data).Error; err != nil {
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

func (db Db) Delete() Db {
	db, tx := db.checkTx()
	if err := tx.First(db.data).Error; err != nil {
		if !db.txMode {
			tx.Rollback()
		}
		db.error = err
		return db
	}

	if err := tx.Delete(db.data).Error; err != nil {
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
