package growl

import (
	"time"

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
	data          interface{}
	where         []dbWhereParams
	join          []dbWhereParams
	or            []dbWhereParams
	selct         string
	preload       []string
	association   string
	limit         int
	offset        int
	orderBy       string
	group         string
	tx            *gorm.DB
	txMode        bool
	growlTag      map[int]map[string]string
	jsonTag       map[int]string
	cacheDuration time.Duration
	error         error
}

type lookUp struct {
	Keys []string
}
