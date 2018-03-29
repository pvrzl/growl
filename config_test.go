package growl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	assert.NotEqual(t, nil, loadConfig("utils_dummy.go"))
	assert.NotEqual(t, nil, loadConfig("utils.go"))
	assert.Equal(t, ErrDbNameRequired, loadConfig("conf.false.yaml"))
	assert.Equal(t, nil, loadConfig("conf.yaml"))
	assert.Equal(t, "mysql", YamlConfig.Growl.Database.Driver)
}

func TestCheckConfig(t *testing.T) {
	assert.Equal(t, nil, loadConfig("conf.yaml"))

	YamlConfig.Growl.Redis.Host = ""
	config, err := checkConfig(YamlConfig)
	assert.Equal(t, nil, err)
	assert.Equal(t, "localhost", config.Growl.Redis.Host)

	YamlConfig.Growl.Redis.Port = ""
	config, err = checkConfig(YamlConfig)
	assert.Equal(t, nil, err)
	assert.Equal(t, "6379", config.Growl.Redis.Port)

	YamlConfig.Growl.Database.Name = ""
	config, err = checkConfig(YamlConfig)
	assert.Equal(t, ErrDbNameRequired, err)
	YamlConfig.Growl.Database.Name = "test"

	YamlConfig.Growl.Database.Driver = ""
	config, err = checkConfig(YamlConfig)
	assert.Equal(t, ErrDbDriverRequired, err)
	YamlConfig.Growl.Database.Driver = "sqlite"

	YamlConfig.Growl.Database.Url = ""
	config, err = checkConfig(YamlConfig)
	assert.Equal(t, ErrDbUrlRequired, err)
}
