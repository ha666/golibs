package golibs

import (
	"sync"
)

var _cache sync.Map

func CacheSet(key, value string) {
	_cache.Store(key, value)
}

func CacheGet(key string) (value string, ok bool) {
	result, ok := _cache.Load(key)
	if !ok {
		return "", false
	}
	return result.(string), true
}

func CacheDel(key string) {
	_cache.Delete(key)
}
