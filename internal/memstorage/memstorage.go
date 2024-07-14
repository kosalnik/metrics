// Package memstorage implements storage in memory.
package memstorage

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/logger"
	"github.com/kosalnik/metrics/internal/storage"

	"github.com/kosalnik/metrics/internal/models"
)

type MemStorage struct {
	updatedAt time.Time
	gauge     map[string]float64
	counter   map[string]int64
	mu        sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:     make(map[string]float64),
		counter:   make(map[string]int64),
		updatedAt: time.Now(),
	}
}

var _ storage.Storage = &MemStorage{}

func (m *MemStorage) GetGauge(_ context.Context, name string) (*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.gauge[name]
	if ok {
		return &models.Metrics{ID: name, MType: models.MGauge, Value: v}, nil
	}

	return nil, nil
}

func (m *MemStorage) GetCounter(_ context.Context, name string) (*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.counter[name]
	if ok {
		return &models.Metrics{ID: name, MType: models.MCounter, Delta: v}, nil
	}

	return nil, nil
}

func (m *MemStorage) SetGauge(_ context.Context, name string, value float64) (*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauge[name] = value
	m.updatedAt = time.Now()

	return &models.Metrics{ID: name, MType: models.MGauge, Value: value}, nil
}

func (m *MemStorage) IncCounter(_ context.Context, name string, value int64) (*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.counter[name] + value
	logger.Logger.WithFields(logrus.Fields{"k": name, "old": m.counter[name], "new": v}).Info("IncCounter")
	m.counter[name] = v
	m.updatedAt = time.Now()

	return &models.Metrics{ID: name, MType: models.MCounter, Delta: v}, nil
}

func (m *MemStorage) UpsertAll(_ context.Context, list []models.Metrics) error {
	if len(list) == 0 {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	logger.Logger.WithField("list", list).Info("upsertAll")
	for _, v := range list {
		switch v.MType {
		case models.MGauge:
			t := v.Value
			m.gauge[v.ID] = t
			m.updatedAt = time.Now()
			continue
		case models.MCounter:
			t := v.Delta
			m.counter[v.ID] += t
			m.updatedAt = time.Now()
		}
	}

	return nil
}

func (m *MemStorage) GetAll(_ context.Context) ([]models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	res := make([]models.Metrics, len(m.gauge)+len(m.counter))
	i := 0
	for k, v := range m.gauge {
		t := v
		res[i] = models.Metrics{ID: k, MType: models.MGauge, Value: t}
		i++
	}
	for k, v := range m.counter {
		t := v
		res[i] = models.Metrics{ID: k, MType: models.MCounter, Delta: t}
		i++
	}
	return res, nil
}

func (m *MemStorage) Close() error {
	return nil
}

func (m *MemStorage) Ping(_ context.Context) error {
	return nil
}

func (m *MemStorage) UpdatedAt() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.updatedAt
}
