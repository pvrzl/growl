package growl

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"
)

func TestDbConnect(t *testing.T) {
	Config.Path = "conf.yaml"
	Config.Load()
	_, err := dbConnect()
	assert.Equal(t, nil, err)

}

func TestDb(t *testing.T) {
	Config.Path = "conf.yaml"
	Config.Load()
	connDb = nil

	_, err := Conn()
	assert.Equal(t, nil, err)

	_, err = Conn()
	assert.Equal(t, nil, err)

	connDb.Close()
	_, err = Conn()
	assert.Equal(t, nil, err)

	connDb.Close()
}

func TestDbSetPrefix(t *testing.T) {
	Config.Path = "conf.yaml"
	Config.Load()
	YamlConfig.Growl.Database.Prefix = "asd_"
	assert.Equal(t, "asd_test", dbSetPrefix(connDb, "test"))
}
