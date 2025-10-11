package context

import (
	"sync"
)

type SafeMap[K comparable, V any] struct {
	data   map[K]V
	locker sync.RWMutex
	keys   *Linked[K]
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
	if _, ok := s.data[key]; !ok {
		if s.keys == nil {
			s.keys = NewLinked(key)
		} else {
			s.keys.Add(key)
		}
	}
	s.data[key] = val
	s.locker.Unlock()
}

func (s *SafeMap[K, V]) Rem(key K) {
	if !s.Exists(key) {
		return
	}

	s.locker.Lock()
	delete(s.data, key)
	s.keys.Rem(key)
	s.locker.Unlock()
}

func (s *SafeMap[K, V]) Exists(key K) bool {
	s.locker.RLock()
	defer s.locker.RUnlock()
	_, ok := s.data[key]
	return ok
}

func (s *SafeMap[K, V]) Len() int {
	s.locker.RLock()
	defer s.locker.RUnlock()

	return len(s.data)
}

func (s *SafeMap[K, V]) Range(f func(key K, val V) bool) {
	s.locker.RLock()
	tmp := make(map[K]V, len(s.data))
	for key, val := range s.data {
		tmp[key] = val
	}
	keys := make([]K, s.keys.Len())
	copy(keys, s.keys.Values())
	s.locker.RUnlock()

	for _, key := range keys {
		if f(key, tmp[key]) {
			break
		}
	}
}
