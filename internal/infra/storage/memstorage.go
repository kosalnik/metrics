package storage

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/kosalnik/metrics/internal/models"
	"github.com/sirupsen/logrus"
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

func NewStorage(backupInterval *time.Duration, backupPath *string) *MemStorage {
	return &MemStorage{
		gauge:          make(map[string]float64),
		counter:        make(map[string]int64),
		backupInterval: backupInterval,
		lastBackup:     time.Now(),
		backupPath:     backupPath,
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
	m.checkBackup()
	return value
}

func (m *MemStorage) IncCounter(name string, value int64) int64 {
	m.mu.Lock()
	v := m.counter[name] + value
	logrus.WithFields(logrus.Fields{"k": name, "old": m.counter[name], "new": v}).Info("IncCounter")
	m.counter[name] = v
	m.mu.Unlock()
	m.checkBackup()
	return v
}

func (m *MemStorage) checkBackup() {
	if m.backupInterval == nil || m.backupPath == nil {
		return
	}
	if *m.backupInterval > 0 && m.lastBackup.Add(*m.backupInterval).Before(time.Now()) {
		return
	}
	if err := m.Store(*m.backupPath); err != nil {
		logrus.WithError(err).Error("failed backup")
	}
	m.lastBackup = time.Now()
}

func (m *MemStorage) GetAll() []models.Metrics {
	m.mu.Lock()
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
	m.mu.Unlock()
	return res
}

type Backup struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (m *MemStorage) Store(path string) error {
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

func (m *MemStorage) Recover(path string) error {
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
