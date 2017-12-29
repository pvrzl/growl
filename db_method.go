package growl

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
