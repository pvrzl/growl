package growl

import (
	"fmt"
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
		// tx.Commit()
	}

	if YamlConfig.Growl.Redis.Enable || YamlConfig.Growl.Misc.LocalCache {
		idv := reflect.ValueOf(db.data).Elem().FieldByName("Id")
		if idv.IsValid() {
			id := valid.ToString(idv.Interface().(int))
			if id != "" {
				DeleteLookup(db.LookupKey(id))
			}
		}
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
			// tx.Rollback()
		}
		return db
	}

	id := reflect.ValueOf(db.data).Elem().FieldByName("Id")
	if id.IsValid() {
		id.Set(reflect.Zero(id.Type()))
	}

	if err := tx.Create(db.data).Error; err != nil {
		if !db.txMode {
			// tx.Rollback()
		}
		db.error = err
		return db
	}

	if !db.txMode {
		// tx.Commit()
	}

	if YamlConfig.Growl.Redis.Enable || YamlConfig.Growl.Misc.LocalCache {
		DeleteLookup(db.LookupKey("count"))
		// DeleteLookup(db.GetTableName())
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
			// tx.Rollback()
		}
		return db
	}

	if err := tx.Save(db.data).Error; err != nil {
		if !db.txMode {
			// tx.Rollback()
		}
		db.error = err
		return db
	}

	if !db.txMode {
		// tx.Commit()
	}

	if YamlConfig.Growl.Redis.Enable || YamlConfig.Growl.Misc.LocalCache {
		idv := reflect.ValueOf(db.data).Elem().FieldByName("Id")
		if idv.IsValid() {
			id := valid.ToString(idv.Interface().(int))
			if id != "" {
				DeleteLookup(db.LookupKey(id))
			}
		}
	}

	return db
}

func (db Db) UpdateMap(data map[string]interface{}) Db {

	db, tx := db.checkTx()

	if err := tx.Update(data).Error; err != nil {
		if !db.txMode {
			// tx.Rollback()
		}
		db.error = err
		return db
	}

	if !db.txMode {
		// tx.Commit()
	}

	if YamlConfig.Growl.Redis.Enable || YamlConfig.Growl.Misc.LocalCache {
		idv := reflect.ValueOf(db.data).Elem().FieldByName("Id")
		if idv.IsValid() {
			id := valid.ToString(idv.Interface().(int))
			if id != "" {
				DeleteLookup(db.LookupKey(id))
			}
		}
	}

	return db
}

func (db Db) Update() Db {
	if _, err := valid.ValidateStruct(db.data); err != nil {
		db.error = err
		return db
	}

	db, tx := db.checkTx()

	db = db.checkTag()
	if db.error != nil {
		if !db.txMode {
			// tx.Rollback()
		}
		return db
	}

	if err := tx.Update(db.data).Error; err != nil {
		if !db.txMode {
			// tx.Rollback()
		}
		db.error = err
		return db
	}

	if !db.txMode {
		// tx.Commit()
	}

	if YamlConfig.Growl.Redis.Enable || YamlConfig.Growl.Misc.LocalCache {
		idv := reflect.ValueOf(db.data).Elem().FieldByName("Id")
		if idv.IsValid() {
			id := valid.ToString(idv.Interface().(int))
			if id != "" {
				DeleteLookup(db.LookupKey(id))
			}
		}
	}

	return db
}

func (db Db) First() Db {
	db.limit = 1
	err := GetCache(MD5(db.GenerateSelectRaw()), db.data)
	fmt.Println("err ", err)
	db, tx := db.checkTx()
	if err == nil {
		if !db.txMode {
			// tx.Commit()
		}
		return db
	}

	if err := tx.First(db.data).Error; err != nil {
		if !db.txMode {
			// tx.Rollback()
		}
		db.error = err
		return db
	}

	if !db.txMode {
		// tx.Commit()
	}

	if YamlConfig.Growl.Redis.Enable || YamlConfig.Growl.Misc.LocalCache {
		key := MD5(db.GenerateSelectRaw())
		SetCache(key, db.data)
		idv := reflect.ValueOf(db.data).Elem().FieldByName("Id")
		if idv.IsValid() {
			id := valid.ToString(idv.Interface().(int))
			if id != "" {
				lu := new(lookUp)
				GetCache(db.LookupKey(id), lu)
				lu.keys = append(lu.keys, key)
				SetCache(db.LookupKey(id), lu)
			}
		}
	}

	return db
}

func (db Db) Find(data interface{}) Db {
	err := GetCache(MD5(db.GenerateSelectRaw()), data)
	db, tx := db.checkTx()
	if err == nil {
		if !db.txMode {
			// tx.Commit()
		}
		return db
	}

	if err := tx.Find(data).Error; err != nil {
		if !db.txMode {
			// tx.Rollback()
		}
		db.error = err
		return db
	}

	if !db.txMode {
		// tx.Commit()
	}

	if YamlConfig.Growl.Redis.Enable || YamlConfig.Growl.Misc.LocalCache {
		key := MD5(db.GenerateSelectRaw())
		SetCache(key, data)
		v := reflect.ValueOf(data).Elem()
		for i := 0; i < v.Len(); i++ {
			idv := v.Index(i).FieldByName("Id")
			if idv.IsValid() {
				id := valid.ToString(idv.Interface().(int))
				if id != "" {
					lu := new(lookUp)
					GetCache(db.LookupKey(id), lu)
					lu.keys = append(lu.keys, key)
					SetCache(db.LookupKey(id), lu)
				}
			}
		}
		// tableLU := new(lookUp)
		// GetCache(db.GetTableName(), tableLU)
		// tableLU.keys = append(tableLU.keys, key)
		// SetCache(db.GetTableName(), tableLU)
	}

	return db
}

func (db Db) Count(data interface{}) Db {
	err := GetCache(MD5(db.GenerateSelectRaw())+"-count", data)
	db, tx := db.checkTx()
	if err == nil {
		if !db.txMode {
			// tx.Commit()
		}
		return db
	}

	if err := tx.Count(data).Error; err != nil {
		if !db.txMode {
			// tx.Rollback()
		}
		db.error = err
		return db
	}

	if !db.txMode {
		// tx.Commit()
	}

	if YamlConfig.Growl.Redis.Enable || YamlConfig.Growl.Misc.LocalCache {
		key := MD5(db.GenerateSelectRaw()) + "-count"
		SetCache(key, data)
		id := db.LookupKey("count")
		lu := new(lookUp)
		GetCache(id, lu)
		lu.keys = append(lu.keys, key)
		SetCache(id, lu)
	}

	return db
}

func (db Db) Delete() Db {
	db, tx := db.checkTx()
	if err := tx.First(db.data).Error; err != nil {
		if !db.txMode {
			// tx.Rollback()
		}
		db.error = err
		return db
	}

	var id string
	idv := reflect.ValueOf(db.data).Elem().FieldByName("Id")
	if idv.IsValid() {
		id = valid.ToString(idv.Interface().(int))
	}

	if err := tx.Delete(db.data).Error; err != nil {
		if !db.txMode {
			// tx.Rollback()
		}
		db.error = err
		return db
	}

	if !db.txMode {
		// tx.Commit()
	}

	if YamlConfig.Growl.Redis.Enable || YamlConfig.Growl.Misc.LocalCache {
		DeleteLookup(db.LookupKey("count"))

		if id != "" {
			DeleteLookup(db.LookupKey(id))
		}

	}

	return db
}
