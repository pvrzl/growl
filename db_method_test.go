package growl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func migrateTestTable() {
	conn, _ := Conn()
	conn.AutoMigrate(TestTable{}, TestTableRelation{})
}

func deleteTestTable() {
	conn, _ := Conn()
	conn.DropTable(TestTable{})
	conn.DropTable(TestTableRelation{})
}

func TestDbWhere(t *testing.T) {
	Config.Load()
	db := Db{}
	db = db.Where("test = ?", "test")
	assert.Equal(t, "test = ?", db.where[0].qry)
	assert.Equal(t, []interface{}{"test"}, db.where[0].params)
	db.Commit()
}

func TestDbSelect(t *testing.T) {
	db := Db{}
	db = db.Select("c1, c2")
	assert.Equal(t, "c1, c2", db.selct)

	db = db.Select("hello world", func(s string) string { return strings.Replace(s, "world", "ef", -1) })
	assert.Equal(t, "hello ef", db.selct)
	db.Commit()
}

func TestDbSave(t *testing.T) {
	Config.Load()
	test := new(TestTable)
	db := test.Db()
	var count int

	// create table
	migrateTestTable()

	db = db.Save()
	assert.NotNil(t, db.error)

	test.Name = "test01"
	db = test.Db().Begin().Save()
	assert.Nil(t, db.error)

	test.Name = "test02"
	test.Db().SetTx(db.GetTx()).Save().Rollback()

	connDb.Model(new(TestTable)).Select("*").Count(&count)
	assert.Equal(t, 0, count)

	test.Name = "test03"
	db = test.Db().Save()
	assert.Nil(t, db.error)
	// fmt.Printf("test data : %+v", test)

	test.Name = "test04"
	db = test.Db().Begin().Save().Commit()

	connDb.Model(new(TestTable)).Select("*").Count(&count)
	assert.Equal(t, 2, count)

	test.Name = "test04"
	db = test.Db().Save()
	connDb.Model(new(TestTable)).Select("*").Count(&count)
	assert.Equal(t, 2, count)
	assert.NotNil(t, db.error)

	testRelation := new(TestTableRelation)
	testRelation.Name = "testRelation01"
	testRelation.TestTableID = 3
	testRelation.Db().Preload("TestTable").Save()
	// fmt.Printf("test relation data : %+v", testRelation)

	connDb.Model(new(TestTableRelation)).Select("*").Count(&count)
	assert.Equal(t, 1, count)
	assert.Equal(t, 3, testRelation.TestTable.Id)
	assert.Equal(t, "test03", testRelation.TestTable.Name)

	assert.Equal(t, 1, connDb.DB().Stats().OpenConnections)

	deleteTestTable()

	// drop table
	// deleteTestTable()
}
