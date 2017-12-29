package growl

var (
	Config     growlConfig
	YamlConfig growlYamlConfig
)

func (config growlConfig) Load() error {
	if config.Path == "" {
		config.Path = "conf.yaml"
		Config = config
	}

	if !IsFileExist(config.Path) {
		return ErrFileNotExist
	}

	return loadConfig(config.Path)

}
