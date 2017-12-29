package growl

type growlConfig struct {
	Path string
}

type dbWhereParams struct {
	qry    string
	params []interface{}
}

type Db struct {
	data  interface{}
	where []dbWhereParams
	selct string
}
