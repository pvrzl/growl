package growl

import (
	"errors"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

func loadConfig(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("error while load " + path + " : " + err.Error())
	}

	err = yaml.Unmarshal(yamlFile, &YamlConfig)
	if err != nil {
		return errors.New("error while unmarshal " + path + " : " + err.Error())
	}

	YamlConfig, err = checkConfig(YamlConfig)
	if err != nil {
		return err
	}

	if YamlConfig.Growl.Misc.FlushAtInit {
		FlushCache()
	}

	return nil
}

func checkConfig(yamlConfig growlYamlConfig) (growlYamlConfig, error) {

	// set database default

	// set redis default
	if yamlConfig.Growl.Redis.Host == "" {
		yamlConfig.Growl.Redis.Host = "localhost"
	}

	if yamlConfig.Growl.Redis.Port == "" {
		yamlConfig.Growl.Redis.Port = "6379"
	}

	yamlConfig.Growl.Redis.Duration = 168 * time.Hour

	if yamlConfig.Growl.Database.Driver == "" {
		return yamlConfig, ErrDbDriverRequired
	}

	if yamlConfig.Growl.Database.Name == "" {
		return yamlConfig, ErrDbNameRequired
	}

	if yamlConfig.Growl.Database.Url == "" {
		return yamlConfig, ErrDbUrlRequired
	}

	return yamlConfig, nil
}
