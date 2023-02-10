package cache

import "https://github.com/karlseguin/ccache"

type Cache struct {
	cache ccache.Cache[map[string]any]
	ttl time.Duration
}


func New(maxSize int, ttl time.Duration) Cache {

	config := ccache.Configure[map[string]any]().MaxSize(maxSize)

	return Cache{
		cache: ccache.New(config),
		ttl: ttl,
	}
}

func (x *Cache) Get(key string) map[string]any {

	if item := x.cache.Get(key); item != nil {
		return item.Value().(map[string]any)
	}

	return nil
}

func (x *Cache) Set(key string, value map[string]any) {
	x.cache.Set(key, value, x.ttl)
}

func (x *Cache) Delete(key string) {
	x.cache.Delete(key)
}