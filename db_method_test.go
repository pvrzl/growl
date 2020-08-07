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
	if db.tx != nil {
		db.Commit()
	}

}

func TestDbSelect(t *testing.T) {
	db := Db{}
	db = db.Select("c1, c2")
	assert.Equal(t, "c1, c2", db.selct)

	db = db.Select("hello world", func(s string) string { return strings.Replace(s, "world", "ef", -1) })
	assert.Equal(t, "hello ef", db.selct)
	if db.tx != nil {
		db.Commit()
	}
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
	testRelation.Db().Save()
	// fmt.Printf("test relation data : %+v", testRelation)

	connDb.Model(new(TestTableRelation)).Select("*").Count(&count)
	assert.Equal(t, 1, count)

	testRelation.TestTableID = 1
	db = testRelation.Db().Save()
	assert.NotNil(t, db.Error())
	// assert.Equal(t, 3, testRelation.TestTable.Id)
	// assert.Equal(t, "test03", testRelation.TestTable.Name)

	assert.Equal(t, 1, OpenConnectionStats())

	deleteTestTable()

	// drop table
	// deleteTestTable()
}

func TestDbFirst(t *testing.T) {
	Config.Load()
	test := new(TestTable)
	var count int

	// create table
	migrateTestTable()

	test.Name = "test01"
	test.Db().Save()

	testRelation := new(TestTableRelation)
	testRelation.Name = "testRelation01"
	testRelation.TestTableID = 1
	testRelation.Db().Save()

	connDb.Model(new(TestTableRelation)).Select("*").Count(&count)
	assert.Equal(t, 1, count)

	testRelation2 := new(TestTableRelation)
	testRelation2.Db().Preload("TestTable").First()

	assert.Equal(t, 1, testRelation2.TestTable.Id)

	db := testRelation2.Db().Preload("TestTable").Where("id = ?", "7").First()
	assert.NotNil(t, db.Error())

	assert.Equal(t, 1, OpenConnectionStats())

	deleteTestTable()

}

func TestDbJoin(t *testing.T) {
	Config.Load()
	migrateTestTable()

	var count int

	test := new(TestTable)
	test.Name = "test01"
	test.Db().Save()

	test.Name = "test02"
	test.Db().Save()

	new(TestTable).Db().Model(new(TestTable)).Select("*").Count(&count)
	assert.Equal(t, 2, count)

	testRelation := new(TestTableRelation)
	testRelation.Name = "testRelation01"
	testRelation.TestTableID = 1
	testRelation.Db().Save()

	new(TestTableRelation).Db().Model(new(TestTableRelation)).Select("*").Count(&count)
	assert.Equal(t, 1, count)

	tests := []TestTable{}

	test.Name = ""
	test.Id = 0
	test.Db().Join("inner join test_table_relations on test_tables.id=test_table_relations.id and test_tables.id = ?", 1).Join("").Where("test_tables.id = ?", 1).Where("test_tables.name = ?", "test01").Find(&tests).Limit(1, func(s int) int { return s }).Offset(0, func(s int) int { return s }).OrderBy("id desc", func(s string) string { return s })
	test.Db().Join("inner join test_table_relations on test_tables.id=test_table_relations.id and test_tables.id = ?", 1).Join("").Where("test_tables.id = ?", 1).Where("test_tables.name = ?", "test01").Find(&tests).Limit(1, func(s int) int { return s }).Offset(0, func(s int) int { return s }).OrderBy("id desc", func(s string) string { return s })

	assert.Equal(t, "test01", tests[0].Name)

	test.Name = "test03"
	test.Db().Save()

	testRelation.Db().Model(new(TestTableRelation)).Association("TestTable").Replace(test)
	testRelation.Db().Preload("TestTable").First()
	assert.Equal(t, "test03", testRelation.TestTable.Name)

	test.Name = ""
	db := test.Db().ForceUpdate()
	assert.NotNil(t, db.error)

	test.Name = "testUpdate"
	test.Db().ForceUpdate()
	test.Db().First()

	assert.Equal(t, "testUpdate", test.Name)

	test.Db().Where("id = ?", test.Id).UpdateMap(map[string]interface{}{"name": "testUpdate02"})
	test.Db().First()
	assert.Equal(t, "testUpdate02", test.Name)

	test.Name = "testUpdate03"
	test.Db().Model(test).Update()
	test.Db().First()
	assert.Equal(t, "testUpdate03", test.Name)

	test.Id = 4
	db = test.Db().Delete()
	// assert.NotNil(t, db.error)

	test.Id = 3
	test.Db().Delete()
	test.Db().Model(test).Count(&count)
	assert.Equal(t, 0, count)

	test.Db().Model(test).Count(&count)
	assert.Equal(t, 0, count)

	deleteTestTable()

}

func TestCache(t *testing.T) {
	Config.Load()
	PingCache()
}

func TestFind(t *testing.T) {
	Config.Load()

	var count int

	migrateTestTable()

	test := new(TestTable)
	test.Name = "test01"
	test.Db().Save()

	test.Name = "test02"
	test.Db().Save()

	test.Db().Model(new(TestTable)).Count(&count)
	assert.Equal(t, 2, count)

	testGet := new(TestTable)
	testGet.Db().Where("id = ?", 1).First()
	assert.Equal(t, "test01", testGet.Name)

	testGet.Id = 0
	testGet.Db().Where("id = ?", 2).Group("").First()
	assert.Equal(t, "test02", testGet.Name)

	deleteTestTable()
}
