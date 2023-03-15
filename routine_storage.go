package routine

import "sync/atomic"

var counter int64

type storage[T interface{}] struct {
	key int64
}

func newStorage[T interface{}]() *storage[T] {
	s := &storage[T]{
		key: atomic.AddInt64(&counter, 1),
	}
	return s
}

func (t *storage[T]) Get() (value T) {
	r := loadStore().get(t.key)
	if r == nil {
		return
	}
	return r.(T)
}

func (t *storage[T]) Set(v T) (oldValue T) {
	r := loadStore().set(t.key, v)
	if r == nil {
		return
	}
	return r.(T)
}

func (t *storage[T]) Del() (oldValue T) {
	r := loadStore().del(t.key)
	if r == nil {
		return
	}
	return r.(T)
}
