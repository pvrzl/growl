package growl

import "time"

type growlYamlConfig struct {
	Growl struct {
		Database struct {
			Driver        string
			Url           string
			Name          string
			Prefix        string
			SingularTable bool
			Limit         int
			TimeLayout    string
		}
		Redis struct {
			Host     string
			Port     string
			Duration time.Duration
			Password string
			Channel  string
		}
		Misc struct {
			LocalCache bool
			Log        bool
		}
	}
}
