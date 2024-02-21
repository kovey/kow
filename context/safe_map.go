package context

import (
	"sync"
)

type SafeMap[K comparable, V any] struct {
	data   map[K]V
	locker sync.RWMutex
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{data: make(map[K]V), locker: sync.RWMutex{}}
}

func (s *SafeMap[K, V]) Get(key K) V {
	var res V
	s.locker.RLock()
	defer s.locker.RUnlock()

	if v, ok := s.data[key]; ok {
		return v
	}

	return res
}

func (s *SafeMap[K, V]) Add(key K, val V) {
	s.locker.Lock()
	s.data[key] = val
	s.locker.Unlock()
}

func (s *SafeMap[K, V]) Rem(key K) {
	s.locker.Lock()
	delete(s.data, key)
	s.locker.Unlock()
}

func (s *SafeMap[K, V]) Len() int {
	s.locker.RLock()
	defer s.locker.RUnlock()

	return len(s.data)
}

func (s *SafeMap[K, V]) For(f func(key K, val V) bool) {
	s.locker.RLock()
	tmp := make(map[K]V, len(s.data))
	for key, val := range s.data {
		tmp[key] = val
	}
	s.locker.RUnlock()

	for key, val := range tmp {
		if f(key, val) {
			break
		}
	}
}
