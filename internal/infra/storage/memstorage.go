package storage

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/models"
)

type MemStorage struct {
	mu             sync.Mutex
	gauge          map[string]float64
	counter        map[string]int64
	backupInterval *time.Duration
	lastBackup     time.Time
	backupPath     *string
}

type MemStorageItem struct {
	class string
	index int
}

func NewMemStorage(backupInterval *time.Duration, backupPath *string) *MemStorage {
	return &MemStorage{
		gauge:          make(map[string]float64),
		counter:        make(map[string]int64),
		backupInterval: backupInterval,
		lastBackup:     time.Now(),
		backupPath:     backupPath,
	}
}

var _ Storage = &MemStorage{}

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
	m.checkBackup(ctx)

	return value, nil
}

func (m *MemStorage) IncCounter(ctx context.Context, name string, value int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v := m.counter[name] + value
	logrus.WithFields(logrus.Fields{"k": name, "old": m.counter[name], "new": v}).Info("IncCounter")
	m.counter[name] = v
	m.checkBackup(ctx)
	return v, nil
}

func (m *MemStorage) UpsertAll(ctx context.Context, list []models.Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	logrus.WithField("list", list).Info("upsertAll")
	for _, v := range list {
		switch v.MType {
		case models.MGauge:
			m.gauge[v.ID] = *v.Value
			continue
		case models.MCounter:
			m.counter[v.ID] += *v.Delta
		}
	}
	m.checkBackup(ctx)
	return nil
}

func (m *MemStorage) checkBackup(ctx context.Context) {
	if m.backupInterval == nil || m.backupPath == nil {
		return
	}
	if *m.backupInterval > 0 && m.lastBackup.Add(*m.backupInterval).Before(time.Now()) {
		return
	}
	if err := m.Store(ctx, *m.backupPath); err != nil {
		logrus.WithError(err).Error("failed backup")
	}
	m.lastBackup = time.Now()
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

type Backup struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (m *MemStorage) Store(_ context.Context, path string) error {
	f, err := os.CreateTemp(os.TempDir(), "backup")
	if err != nil {
		return err
	}
	savePath := f.Name()
	d, err := json.Marshal(Backup{Gauge: m.gauge, Counter: m.counter})
	if err != nil {
		return err
	}
	if _, err := f.Write(d); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Rename(savePath, path); err != nil {
		return err
	}
	return nil
}

func (m *MemStorage) Recover(_ context.Context, path string) error {
	d, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var b Backup
	if err := json.Unmarshal(d, &b); err != nil {
		return err
	}
	m.gauge = b.Gauge
	m.counter = b.Counter
	return nil
}
