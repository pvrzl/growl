package growl

import (
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"

	gocache "github.com/patrickmn/go-cache"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

var connRedis *redis.Client
var localCache = gocache.New(24*time.Hour, 30*time.Minute)

func connectRedis() *redis.Client {
	config := YamlConfig.Growl
	return redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       0,
	})
}

func Redis() *redis.Client {
	if connRedis == nil {
		connRedis = connectRedis()
		return connRedis
	}

	_, err := connRedis.Ping().Result()
	if err != nil {
		connRedis.Close()
		connRedis = connectRedis()
	}

	return connRedis

}

func Codec() *cache.Codec {
	return &cache.Codec{
		Redis: Redis(),
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

func PingCache() error {
	_, err := Redis().Ping().Result()
	return err
}

func FlushCache() {
	Redis().FlushAll()
}

func GetCache(key string, data interface{}) (err error) {
	config := YamlConfig.Growl

	if config.Redis.Enable {
		err = Codec().Get(key, data)
	}

	if !config.Redis.Enable {
		err = ErrCacheDisabled
	}

	return
}

func SetCache(key string, data interface{}) {
	config := YamlConfig.Growl

	if config.Redis.Enable {
		Codec().Set(&cache.Item{
			Key:        key,
			Object:     data,
			Expiration: config.Redis.Duration,
		})
	}
}

func DeleteCache(key string) {
	config := YamlConfig.Growl
	if config.Redis.Enable {
		Codec().Delete(key)
	}
}
