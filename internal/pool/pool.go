package pool

import "sync"

type Resettable interface {
	Reset()
}

type Pool[T any] struct {
	pool  *sync.Pool
	reset func(t T)
}

func NewResettablePool[T Resettable](ctor func() T) *Pool[T] {
	return &Pool[T]{
		pool: newSyncPool(ctor),
		reset: func(t T) {
			t.Reset()
		},
	}
}

func NewPool[T any](ctor func() T) *Pool[T] {
	return &Pool[T]{
		pool:  newSyncPool(ctor),
		reset: func(t T) {},
	}
}

func newSyncPool[T any](ctor func() T) *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return ctor()
		},
	}
}

func (p *Pool[T]) Get() T {
	return p.pool.Get().(T) //nolint: errcheck // always correct type
}

func (p *Pool[T]) Put(t T) {
	p.reset(t)
	p.pool.Put(t)
}
