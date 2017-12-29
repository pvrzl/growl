package growl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTableName(t *testing.T) {
	test := new(TestTable)

	assert.Equal(t, "test_tables", test.Db().GetTableName())
	YamlConfig.Growl.Database.SingularTable = true

	assert.Equal(t, "test_table", test.Db().GetTableName())
}
