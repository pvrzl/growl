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

	connRedis = Redis()
	codec = Codec()

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

	if yamlConfig.Growl.Misc.DefaultCacheDuration != 0 && yamlConfig.Growl.Misc.DefaultCacheDurationUnit != "" {
		switch yamlConfig.Growl.Misc.DefaultCacheDurationUnit {
		case "hour":
			yamlConfig.Growl.Redis.duration = time.Duration(yamlConfig.Growl.Misc.DefaultCacheDuration) * time.Hour
		case "minute":

			yamlConfig.Growl.Redis.duration = time.Duration(yamlConfig.Growl.Misc.DefaultCacheDuration) * time.Minute
		case "second":

			yamlConfig.Growl.Redis.duration = time.Duration(yamlConfig.Growl.Misc.DefaultCacheDuration) * time.Second
		default:

			yamlConfig.Growl.Redis.duration = time.Duration(yamlConfig.Growl.Misc.DefaultCacheDuration) * time.Minute
		}
	} else {
		yamlConfig.Growl.Redis.duration = 1 * time.Hour
	}

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
