package growl

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTableName(t *testing.T) {
	test := new(TestTable)

	assert.Equal(t, "test_tables", test.Db().GetTableName())
	YamlConfig.Growl.Database.SingularTable = true

	assert.Equal(t, "test_table", test.Db().GetTableName())
}

func TestGetError(t *testing.T) {
	test := new(TestTable)
	db := test.Db()
	assert.Equal(t, nil, db.Error())

	testError := errors.New("error")
	db.error = testError
	assert.Equal(t, testError, db.Error())
}

func TestTx(t *testing.T) {

	test := new(TestTable)

	db := test.Db().SetData(test)
	db = db.Begin()
	assert.NotEqual(t, nil, db.tx)
	assert.Equal(t, true, db.txMode)

	db = db.SetTx(db.GetTx())
	assert.NotEqual(t, nil, db.tx)
	assert.Equal(t, true, db.txMode)

	db = db.SetTx(nil)
	assert.Nil(t, db.tx)
	assert.Equal(t, false, db.txMode)

	assert.Nil(t, db.GetTx())
}

func TestGetValue(t *testing.T) {
	a := "test"
	assert.Equal(t, "test", GetValue(reflect.ValueOf(&a)).(string))

	var b string
	assert.Equal(t, "", GetValue(reflect.ValueOf(&b)).(string))
}
