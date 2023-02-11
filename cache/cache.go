package cache

import (
	"time"

	"github.com/karlseguin/ccache/v3"
)

type Cache struct {
	cache *ccache.Cache[map[string]any]
	ttl   time.Duration
}

func New(maxSize int64, ttl time.Duration) *Cache {

	config := ccache.Configure[map[string]any]().MaxSize(maxSize)

	return &Cache{
		cache: ccache.New(config),
		ttl:   ttl,
	}
}

func NewDefaultCache() *Cache {
	return New(2048, 8*time.Minute)
}

func (x *Cache) Get(key string) map[string]any {

	if item := x.cache.Get(key); item != nil {
		return item.Value()
	}

	return nil
}

func (x *Cache) Set(key string, value map[string]any) {
	x.cache.Set(key, value, x.ttl)
}

func (x *Cache) Delete(key string) {
	x.cache.Delete(key)
}
