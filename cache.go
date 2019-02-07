package growl

import (
	"reflect"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"

	gocache "github.com/patrickmn/go-cache"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

var connRedis *redis.Client
var LocalCache = gocache.New(YamlConfig.Growl.Redis.duration, 30*time.Minute)

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
	config := YamlConfig.Growl

	if config.Misc.LocalCache {
		LocalCache.Flush()
	}

	if config.Redis.Enable {
		Redis().FlushDB()
	}
}

func GetCache(key string, data interface{}) (err error) {
	config := YamlConfig.Growl

	if config.Misc.LocalCache {
		cacheData, found := LocalCache.Get(key)
		if config.Misc.Log {
			// fmt.Println("get local key", key)
			// log.Println("get local ", key, " found : ", found)
			// fmt.Println("get local cachedata", cacheData)
		}

		if found {
			x := reflect.ValueOf(data)
			x.Elem().Set(reflect.ValueOf(cacheData).Elem())
			return
		} else {
			err = ErrCacheNotFound
		}
	}

	if config.Redis.Enable {
		err = Codec().Get(key, data)
		if config.Misc.Log {
			// fmt.Println("get redis key", key)
			// log.Println("get redis ", key, " error : ", err)
			// fmt.Println("get redis data", data)
		}

		if err == nil {
			LocalCache.Set(key, data, YamlConfig.Growl.Redis.duration)
			return
		}
	}

	if !config.Redis.Enable && !config.Misc.LocalCache {
		err = ErrCacheDisabled
	}

	return
}

func SetCache(key string, data interface{}, options ...interface{}) {
	config := YamlConfig.Growl

	duration := config.Redis.duration

	if len(options) >= 1 && len(options) == 1 {
		if val, ok := options[0].(time.Duration); ok {
			if val.String() != "0s" {
				duration = val
			}
		}
	}

	if config.Misc.LocalCache {
		LocalCache.Set(key, data, duration)
	}

	if config.Redis.Enable {
		Codec().Set(&cache.Item{
			Key:        key,
			Object:     data,
			Expiration: duration,
		})
	}
}

func DeleteCache(key string) {
	config := YamlConfig.Growl
	if config.Misc.LocalCache {
		LocalCache.Delete(key)
	}
	if config.Redis.Enable {
		Codec().Delete(key)
	}
}
