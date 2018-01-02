package golibs

import (
	"sync"
	"time"
)

var _cache sync.Map

func CacheSet(key string, value interface{}) {
	_cache.Store(key, value)
}

func CacheGetInt(key string) (value int, ok bool) {
	result, ok := _cache.Load(key)
	if !ok {
		return -1, false
	}
	return result.(int), true
}

func CacheGetInt64(key string) (value int64, ok bool) {
	result, ok := _cache.Load(key)
	if !ok {
		return -1, false
	}
	return result.(int64), true
}

func CacheGetString(key string) (value string, ok bool) {
	result, ok := _cache.Load(key)
	if !ok {
		return "", false
	}
	return result.(string), true
}

func CacheGetBool(key string) (value bool, ok bool) {
	result, ok := _cache.Load(key)
	if !ok {
		return false, false
	}
	return result.(bool), true
}

func CacheGetTime(key string) (value time.Time, ok bool) {
	result, ok := _cache.Load(key)
	if !ok {
		return time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), false
	}
	return result.(time.Time), true
}

func CacheGetByteArray(key string) (value []byte, ok bool) {
	result, ok := _cache.Load(key)
	if !ok {
		return nil, false
	}
	return result.([]byte), true
}

func CacheDel(key string) {
	_cache.Delete(key)
}

func CacheCount() int {
	count := 0
	_cache.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}
