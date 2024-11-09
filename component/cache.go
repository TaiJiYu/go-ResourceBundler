package component

import (
	"sync"
	"time"

	"github.com/allegro/bigcache"
)

var cacheTimeout = time.Minute * 10

type cacheClient struct {
	cache *bigcache.BigCache
}

var cacheCli *cacheClient
var cacheOnce sync.Once

func defaultCache() *cacheClient {
	cacheOnce.Do(func() {
		cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(cacheTimeout))
		if err != nil {
			panic(err)
		}
		cacheCli = &cacheClient{
			cache: cache,
		}
	})
	return cacheCli
}

func (c *cacheClient) get(key string) ([]byte, bool) {
	v, err := c.cache.Get(key)
	return v, err == nil
}

func (c *cacheClient) set(key string, value []byte) {
	c.cache.Set(key, value)
}
