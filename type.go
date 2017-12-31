package growl

import (
	"github.com/jinzhu/gorm"
)

type growlConfig struct {
	Path string
}

type dbWhereParams struct {
	qry    string
	params []interface{}
}

type Db struct {
	data    interface{}
	where   []dbWhereParams
	selct   string
	preload []string
	limit   int
	tx      *gorm.DB
	txMode  bool
	error   error
}
