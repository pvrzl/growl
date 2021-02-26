package growl

import (
	"reflect"
	"sync"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"

	gocache "github.com/patrickmn/go-cache"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

var connRedis *redis.Client
var LocalCache = gocache.New(YamlConfig.Growl.Redis.duration, 30*time.Minute)
var redisOnce sync.Once

type codecStruct struct {
	codec *cache.Codec
	sync  sync.Mutex
}

var codec *codecStruct

func Cache() *codecStruct {
	return codec
}

func connectRedis() *redis.Client {
	config := YamlConfig.Growl
	return redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       0,
	})
}

func Redis() *redis.Client {

	redisOnce.Do(func() {
		connRedis = connectRedis()
	})

	return connRedis

}

func Codec() *codecStruct {
	return &codecStruct{
		codec: &cache.Codec{
			Redis: connRedis,
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)
			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
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

func GetCache(key string, data interface{}, options ...interface{}) error {
	return codec.GetCache(key, data)
}

func (codec *codecStruct) GetCache(key string, data interface{}) (err error) {
	codec.sync.Lock()
	defer codec.sync.Unlock()
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
		err = codec.codec.Get(key, data)
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
	codec.SetCache(key, data, options...)
}

func (codec *codecStruct) SetCache(key string, data interface{}, options ...interface{}) {
	codec.sync.Lock()
	defer codec.sync.Unlock()
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
		codec.codec.Set(&cache.Item{
			Key:        key,
			Object:     data,
			Expiration: duration,
		})
	}
}

func (codec *codecStruct) DeleteCache(key string) {
	codec.sync.Lock()
	defer codec.sync.Unlock()
	config := YamlConfig.Growl
	if config.Misc.LocalCache {
		LocalCache.Delete(key)
	}
	if config.Redis.Enable {
		codec.codec.Delete(key)
	}
}
