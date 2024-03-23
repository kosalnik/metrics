package storage

import (
	"sync"

	"github.com/kosalnik/metrics/internal/models"
	"github.com/sirupsen/logrus"
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
	logrus.WithFields(logrus.Fields{"k": name, "old": m.counter[name], "new": v}).Info("IncCounter")
	m.counter[name] = v
	m.mu.Unlock()
	return v
}

func (m *MemStorage) GetAll() []models.Metrics {
	m.mu.Lock()
	res := make([]models.Metrics, len(m.gauge)+len(m.counter))
	i := 0
	for k, v := range m.gauge {
		t := v
		res[i] = models.Metrics{ID: k, MType: "gauge", Value: &t}
		i++
	}
	for k, v := range m.counter {
		t := v
		res[i] = models.Metrics{ID: k, MType: "counter", Delta: &t}
		i++
	}
	m.mu.Unlock()
	return res
}
