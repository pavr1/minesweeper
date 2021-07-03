package ccache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CCache struct {
	Cache *cache.Cache
}

func New() *CCache {
	cache := &CCache{
		Cache: cache.New(1*time.Hour, 60*time.Minute),
	}

	return cache
}

func (c *CCache) Set(key string, data interface{}) {
	c.Cache.Set(key, data, cache.NoExpiration)
}

func (c *CCache) Get(key string) (interface{}, bool) {
	return c.Cache.Get(key)
}

func (c *CCache) Flush() {
	c.Cache.Flush()
}
