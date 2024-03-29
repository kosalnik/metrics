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

func (m *MemStorage) GetGauge(name string) (v float64, ok bool) {
	v, ok = m.gauge[name]
	return
}

func (m *MemStorage) GetCounter(name string) (v int64, ok bool) {
	v, ok = m.counter[name]
	return
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
