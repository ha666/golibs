package golibs

import (
	"sync"
)

var _cache sync.Map

func CacheSet(key, value string) {
	_cache.Store(key, value)
}

func CacheGet(key string) (value string, errmsg string) {
	result, ok := _cache.Load(key)
	if !ok {
		return "", "not found"
	}
	return result.(string), ""
}

func CacheDel(key string) {
	_cache.Delete(key)
}
