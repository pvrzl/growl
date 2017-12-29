package growl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFileExist(t *testing.T) {
	assert.Equal(t, true, IsFileExist("utils.go"))
	assert.Equal(t, false, IsFileExist("utils_dummy.go"))
}

func TestGetName(t *testing.T) {
	test01 := TestTable{}
	assert.Equal(t, "TestTable", GetStructName(test01))
	test02 := new(TestTable)
	assert.Equal(t, "TestTable", GetStructName(test02))
}
