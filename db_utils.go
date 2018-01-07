package growl

import (
	"errors"
	"reflect"
	"strings"

	valid "github.com/asaskevich/govalidator"
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

func (db Db) checkTag() Db {
	v := reflect.ValueOf(db.data).Elem()
	t := reflect.TypeOf(db.data).Elem()

	db = db.getTag(v, t)
	for i := 0; i < v.NumField(); i++ {
		if growl, ok := db.growlTag[i]; ok {
			if growlValue, ok2 := growl["exist"]; ok2 {
				value := GetValue(v.Field(i))
				var dummy struct{}
				_, tx := db.checkTx()
				err := tx.Table(growlValue).Select(growl["existColumn"]).Where(growl["existColumn"]+" = ?", value).First(&dummy).Error
				if err != nil {
					if !db.txMode {
						tx.Rollback()
					}
					db.error = errors.New("error on processing " + t.Field(i).Name + " : " + err.Error())
					return db
				}
			}
		}
	}

	return db
}

func (db Db) getTag(v reflect.Value, t reflect.Type) Db {

	growlTag := make(map[int]map[string]string)
	jsonTag := make(map[int]string)

	for i := 0; i < v.NumField(); i++ {
		growls, ok := t.Field(i).Tag.Lookup("growl")
		if ok {
			growlBody := make(map[string]string)
			for _, growl := range strings.Split(growls, ";") {
				kv := strings.Split(growl, ":")
				growlBody[kv[0]] = kv[1]
			}
			growlTag[i] = growlBody
		}
		json, ok := t.Field(i).Tag.Lookup("json")
		if ok {
			jsonTag[i] = json
		}
	}

	db.growlTag = growlTag
	db.jsonTag = jsonTag
	return db
}

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

func (db Db) Model(data interface{}) Db {
	db, tx := db.checkTx()
	db.tx = tx.Model(data)
	return db
}

func (db Db) GenerateSelectRaw() string {
	table := db.GetTableName()
	var where, join, selct, offset, limit, order, raw string

	for i, obj := range db.where {
		if i != 0 {
			where += " AND "
		}
		where += "(" + obj.qry + ")"
		for _, param := range obj.params {
			strings.Replace(where, "?", valid.ToString(param), 1)
		}
	}

	for i, obj := range db.join {
		if i != 0 {
			join += " AND "
		}
		join += "(" + obj.qry + ")"
		for _, param := range obj.params {
			strings.Replace(join, "?", valid.ToString(param), 1)
		}
	}

	if db.selct == "" {
		selct = "*"
	} else {
		selct = db.selct
	}

	offset = " OFFSET " + valid.ToString(db.offset)
	limit = " LIMIT " + valid.ToString(db.limit)
	order = db.orderBy

	raw = "[ SELECT " + selct + " FROM " + table + join + " WHERE " + where + limit + offset + order + " ][ Preload : " + strings.Join(db.preload, ",") + " ]"
	return raw
}
