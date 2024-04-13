package memstorage

import (
	"context"
	"sync"
	"time"

	"github.com/kosalnik/metrics/internal/infra/logger"
	"github.com/kosalnik/metrics/internal/infra/storage"
	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/models"
)

type MemStorage struct {
	mu        sync.Mutex
	gauge     map[string]float64
	counter   map[string]int64
	updatedAt time.Time
}

type MemStorageItem struct {
	class string
	index int
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:     make(map[string]float64),
		counter:   make(map[string]int64),
		updatedAt: time.Now(),
	}
}

var _ storage.Storage = &MemStorage{}

func (m *MemStorage) GetGauge(_ context.Context, name string) (v float64, ok bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok = m.gauge[name]

	return
}

func (m *MemStorage) GetCounter(_ context.Context, name string) (v int64, ok bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok = m.counter[name]

	return
}

func (m *MemStorage) SetGauge(ctx context.Context, name string, value float64) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauge[name] = value
	m.updatedAt = time.Now()

	return value, nil
}

func (m *MemStorage) IncCounter(ctx context.Context, name string, value int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.counter[name] + value
	logger.Logger.WithFields(logrus.Fields{"k": name, "old": m.counter[name], "new": v}).Info("IncCounter")
	m.counter[name] = v
	m.updatedAt = time.Now()

	return v, nil
}

func (m *MemStorage) UpsertAll(ctx context.Context, list []models.Metrics) error {
	if len(list) == 0 {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	logger.Logger.WithField("list", list).Info("upsertAll")
	for _, v := range list {
		switch v.MType {
		case models.MGauge:
			m.gauge[v.ID] = *v.Value
			m.updatedAt = time.Now()
			continue
		case models.MCounter:
			m.counter[v.ID] += *v.Delta
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
		res[i] = models.Metrics{ID: k, MType: models.MGauge, Value: &t}
		i++
	}
	for k, v := range m.counter {
		t := v
		res[i] = models.Metrics{ID: k, MType: models.MCounter, Delta: &t}
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
