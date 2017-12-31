package growl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestTable struct {
	Name string `valid:"required" gorm:"unique_index"`
	Id   int    `gorm:"AUTO_INCREMENT"`
}

func (test *TestTable) Db() (db Db) {
	db.data = test
	return db
}

type TestTableRelation struct {
	Name        string    `valid:"required"`
	Id          int       `gorm:"AUTO_INCREMENT" json:"id"`
	TestTable   TestTable `gorm:"ForeignKey:TestTableID;AssociationForeignKey:ID" valid:"-"`
	TestTableID int       `valid:"required" growl:"exist:test_tables;existColumn:id"`
}

func (test *TestTableRelation) Db() (db Db) {
	db.data = test
	return db
}

func TestLoad(t *testing.T) {
	Config.Path = "utils_dummy.go"
	assert.Equal(t, ErrFileNotExist, Config.Load())

	Config.Path = "utils.go"
	assert.NotEqual(t, nil, Config.Load())

	Config.Path = ""
	assert.Equal(t, nil, Config.Load())

}
