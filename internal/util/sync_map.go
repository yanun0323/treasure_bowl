package util

import "sync"

type SyncMap[K comparable, V any] struct {
	lock *sync.RWMutex
	data map[K]V
}

func NewSyncMap[K comparable, V any]() SyncMap[K, V] {
	return SyncMap[K, V]{
		lock: &sync.RWMutex{},
		data: map[K]V{},
	}
}

func (m *SyncMap[K, V]) Load(key K) (V, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	d, ok := m.data[key]
	return d, ok
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.data[key] = value
}

func (m *SyncMap[K, V]) LoadAndSet(key K, fn func(value V) V) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.data[key] = fn(m.data[key])
}

func (m *SyncMap[K, V]) Iter(fn func(key K, value V) bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for k, v := range m.data {
		if !fn(k, v) {
			return
		}
	}
}

func (m *SyncMap[K, V]) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.data)
}

func (m *SyncMap[K, V]) Clone() map[K]V {
	m.lock.RLock()
	defer m.lock.RUnlock()
	c := make(map[K]V, len(m.data))
	for k, v := range m.data {
		c[k] = v
	}
	return c
}

func (m *SyncMap[K, V]) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	clear(m.data)
}
