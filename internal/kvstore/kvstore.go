package kvstore

import (
	"sync"
)

type Store[K comparable, V any] struct {
	lock sync.RWMutex

	m map[K]V
}

func NewStore[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{
		m: make(map[K]V),
	}
}

func (s *Store[K, V]) Rand() (K, V, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for k, v := range s.m {
		return k, v, true
	}

	var k K
	var v V

	return k, v, false
}

func (d *Store[K, V]) Set(k K, v V) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.m[k] = v
}

func (d *Store[K, V]) Get(k K) (V, bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	v, ok := d.m[k]

	return v, ok
}

func (d *Store[K, V]) Del(k K) (V, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()

	v, ok := d.m[k]

	delete(d.m, k)

	return v, ok

}
