package growl

import "errors"

var (
	ErrFileNotExist     = errors.New("config file not exist(default : conf.yaml)")
	ErrDbDriverRequired = errors.New("database driver is required")
	ErrDbUrlRequired    = errors.New("database url is required")
	ErrDbNameRequired   = errors.New("database name is required")
	ErrCacheDisabled    = errors.New("cache is disabled")
	ErrCacheNotFound    = errors.New("cache not found")
)
