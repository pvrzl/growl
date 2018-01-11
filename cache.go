package growl

import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"

	gocache "github.com/patrickmn/go-cache"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

var connRedis *redis.Client
var LocalCache = gocache.New(24*time.Hour, 30*time.Minute)

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
	LocalCache.Flush()
}

func GetCache(key string, data interface{}) (err error) {
	config := YamlConfig.Growl

	if config.Misc.LocalCache {
		fmt.Println("get key", key)
		cacheData, found := LocalCache.Get(key)
		fmt.Println("get found", found)
		fmt.Println("get cachedata", cacheData)
		if !found {
			err = ErrCacheNotFound
			return
		} else {
			x := reflect.ValueOf(data)
			x.Elem().Set(reflect.ValueOf(cacheData).Elem())
			return
		}
	}

	if config.Redis.Enable {
		err = Codec().Get(key, data)
		if err == nil {
			LocalCache.Set(key, data, gocache.DefaultExpiration)
			return
		}
	}

	if !config.Redis.Enable && !config.Misc.LocalCache {
		err = ErrCacheDisabled
	}

	return
}

func SetCache(key string, data interface{}) {
	fmt.Println("set key", key)
	config := YamlConfig.Growl

	if config.Misc.LocalCache {
		LocalCache.Set(key, data, gocache.DefaultExpiration)
	}

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
	if config.Misc.LocalCache {
		LocalCache.Delete(key)
	}
	if config.Redis.Enable {
		Codec().Delete(key)
	}
}
