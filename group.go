package golibs

import (
	"context"
	"math"
	"sync"
)

type LimitWaitGroup struct {
	size    int
	current chan struct{}
	wg      sync.WaitGroup
}

func NewLimitWaitGroup(limit int) *LimitWaitGroup {
	size := math.MaxInt32 // 2^32 - 1
	if limit > 0 {
		size = limit
	}

	return &LimitWaitGroup{
		size:    size,
		current: make(chan struct{}, size),
		wg:      sync.WaitGroup{},
	}
}

func (s *LimitWaitGroup) Add() error {
	return s.AddWithContext(context.Background())
}

func (s *LimitWaitGroup) AddWithContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.current <- struct{}{}:
		break
	}
	s.wg.Add(1)
	return nil
}

func (s *LimitWaitGroup) Done() {
	<-s.current
	s.wg.Done()
}

func (s *LimitWaitGroup) Wait() {
	s.wg.Wait()
}
