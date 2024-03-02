package storage

import (
	"fmt"
)

type MemStorage struct {
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

func (m *MemStorage) HasGauge(name string) bool {
	_, ok := m.gauge[name]
	return ok
}

func (m *MemStorage) GetGauge(name string) float64 {
	return m.gauge[name]
}

func (m *MemStorage) HasCounter(name string) bool {
	_, ok := m.counter[name]
	return ok
}

func (m *MemStorage) GetCounter(name string) int64 {
	return m.counter[name]
}

func (m *MemStorage) SetGauge(name string, value float64) {
	m.gauge[name] = value
}

func (m *MemStorage) IncCounter(name string, value int64) {
	m.counter[name] = m.counter[name] + value
}

func (m *MemStorage) GetPlain() map[string]string {
	res := make(map[string]string, len(m.gauge)+len(m.counter))
	for k, v := range m.gauge {
		res[k] = fmt.Sprintf("%v", v)
	}
	for k, v := range m.counter {
		res[k] = fmt.Sprintf("%v", v)
	}
	return res
}
