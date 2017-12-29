package growl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestTable struct {
	Name string
}

func (test *TestTable) Db() (db Db) {
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
