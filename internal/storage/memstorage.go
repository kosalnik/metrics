package storage

import (
	"fmt"
	"sync"
)

type MemStorage struct {
	mu      sync.Mutex
	gauge   map[string]float64
	counter map[string]int64
}

type MemStorageItem struct {
	class string
	index int
}

func NewStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (m *MemStorage) GetGauge(name string) (v float64, ok bool) {
	m.mu.Lock()
	v, ok = m.gauge[name]
	m.mu.Unlock()
	return
}

func (m *MemStorage) GetCounter(name string) (v int64, ok bool) {
	m.mu.Lock()
	v, ok = m.counter[name]
	m.mu.Unlock()
	return
}

func (m *MemStorage) SetGauge(name string, value float64) float64 {
	m.mu.Lock()
	m.gauge[name] = value
	m.mu.Unlock()
	return value
}

func (m *MemStorage) IncCounter(name string, value int64) int64 {
	m.mu.Lock()
	v := m.counter[name] + value
	m.counter[name] = v
	m.mu.Unlock()
	return v
}

func (m *MemStorage) GetPlain() map[string]string {
	m.mu.Lock()
	res := make(map[string]string, len(m.gauge)+len(m.counter))
	for k, v := range m.gauge {
		res[k] = fmt.Sprintf("%v", v)
	}
	for k, v := range m.counter {
		res[k] = fmt.Sprintf("%v", v)
	}
	m.mu.Unlock()
	return res
}
