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
		}
		Redis struct {
			Host     string
			Port     string
			Password string
			Channel  string
			Enable   bool
			duration time.Duration
		}
		Misc struct {
			LocalCache               bool
			Log                      bool
			FlushAtInit              bool
			Debug                    bool
			DefaultCacheDuration     uint
			DefaultCacheDurationUnit string
		}
	}
}
