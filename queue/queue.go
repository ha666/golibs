package queue

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

type esCache struct {
	putNo int64
	getNo int64
	value interface{}
}

// lock free queue
type EsQueue struct {
	capaciity int64
	capMod    int64
	putPos    int64
	getPos    int64
	cache     []esCache
}

func NewQueue(capaciity int64) *EsQueue {
	q := new(EsQueue)
	q.capaciity = minQuantity(capaciity)
	q.capMod = q.capaciity - 1
	q.putPos = 0
	q.getPos = 0
	q.cache = make([]esCache, q.capaciity)
	for i := range q.cache {
		cache := &q.cache[i]
		cache.getNo = int64(i)
		cache.putNo = int64(i)
	}
	cache := &q.cache[0]
	cache.getNo = q.capaciity
	cache.putNo = q.capaciity
	return q
}

func (q *EsQueue) String() string {
	getPos := atomic.LoadInt64(&q.getPos)
	putPos := atomic.LoadInt64(&q.putPos)
	return fmt.Sprintf("Queue{capaciity: %v, capMod: %v, putPos: %v, getPos: %v}",
		q.capaciity, q.capMod, putPos, getPos)
}

func (q *EsQueue) Capaciity() int64 {
	return q.capaciity
}

func (q *EsQueue) Len() int64 {
	var putPos, getPos int64
	var quantity int64
	getPos = atomic.LoadInt64(&q.getPos)
	putPos = atomic.LoadInt64(&q.putPos)

	if putPos >= getPos {
		quantity = putPos - getPos
	} else {
		quantity = q.capMod + (putPos - getPos)
	}

	return quantity
}

// put queue functions
func (q *EsQueue) Push(val interface{}) (ok bool, quantity int64) {
	var putPos, putPosNew, getPos, posCnt int64
	var cache *esCache
	capMod := q.capMod

	getPos = atomic.LoadInt64(&q.getPos)
	putPos = atomic.LoadInt64(&q.putPos)

	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	if posCnt >= capMod-1 {
		runtime.Gosched()
		return false, posCnt
	}

	putPosNew = putPos + 1
	if !atomic.CompareAndSwapInt64(&q.putPos, putPos, putPosNew) {
		runtime.Gosched()
		return false, posCnt
	}

	cache = &q.cache[putPosNew&capMod]

	for {
		getNo := atomic.LoadInt64(&cache.getNo)
		putNo := atomic.LoadInt64(&cache.putNo)
		if putPosNew == putNo && getNo == putNo {
			cache.value = val
			atomic.AddInt64(&cache.putNo, q.capaciity)
			return true, posCnt + 1
		} else {
			runtime.Gosched()
		}
	}
}

// get queue functions
func (q *EsQueue) Pop() (val interface{}, ok bool, quantity int64) {
	var putPos, getPos, getPosNew, posCnt int64
	var cache *esCache
	capMod := q.capMod

	putPos = atomic.LoadInt64(&q.putPos)
	getPos = atomic.LoadInt64(&q.getPos)

	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	if posCnt < 1 {
		runtime.Gosched()
		return nil, false, posCnt
	}

	getPosNew = getPos + 1
	if !atomic.CompareAndSwapInt64(&q.getPos, getPos, getPosNew) {
		runtime.Gosched()
		return nil, false, posCnt
	}

	cache = &q.cache[getPosNew&capMod]

	for {
		getNo := atomic.LoadInt64(&cache.getNo)
		putNo := atomic.LoadInt64(&cache.putNo)
		if getPosNew == getNo && getNo == putNo-q.capaciity {
			val = cache.value
			atomic.AddInt64(&cache.getNo, q.capaciity)
			return val, true, posCnt - 1
		} else {
			runtime.Gosched()
		}
	}
}

// round 到最近的2的倍数
func minQuantity(v int64) int64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
