package growl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func migrateTestTable() {
	conn, _ := Conn()
	conn.AutoMigrate(TestTable{})
}

func deleteTestTable() {
	conn, _ := Conn()
	conn.DropTable(TestTable{})
}

func TestDbWhere(t *testing.T) {
	db := Db{}
	db = db.Where("test = ?", "test")
	assert.Equal(t, "test = ?", db.where[0].qry)
	assert.Equal(t, []interface{}{"test"}, db.where[0].params)
}

func TestDbSelect(t *testing.T) {
	db := Db{}
	db = db.Select("c1, c2")
	assert.Equal(t, "c1, c2", db.selct)
}

func TestDbSave(t *testing.T) {
	Config.Load()
	migrateTestTable()
	deleteTestTable()
}
