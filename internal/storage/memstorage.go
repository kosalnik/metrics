package storage

import (
	"github.com/kosalnik/metrics/internal/entity"
	"log"
)

type Storage interface {
	GetGauge(name string) *entity.GaugeValue
	GetCounter(name string) *entity.CounterValue
	SetGauge(name string, value float64)
	IncCounter(name string, value int64)
}

type MemStorage struct {
	gauge   map[string]*entity.GaugeValue
	counter map[string]*entity.CounterValue
}

type MemStorageItem struct {
	class string
	index int
}

var s = NewStorage()

func NewStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]*entity.GaugeValue),
		counter: make(map[string]*entity.CounterValue),
	}
}

func GetStorage() *MemStorage {
	return s
}

func (m *MemStorage) GetGauge(name string) *entity.GaugeValue {
	if ref, ok := m.gauge[name]; ok {
		return ref
	}
	return nil
}

func (m *MemStorage) GetCounter(name string) *entity.CounterValue {
	if ref, ok := m.counter[name]; ok {
		return ref
	}
	return nil
}

func (m *MemStorage) SetGauge(name string, value float64) {
	if item := m.GetGauge(name); item != nil {
		item.Value = value
		log.Printf("SetGauge[%s]=%v\n", name, m.gauge[name].Value)
		return
	}
	item := entity.GaugeValue{Name: name, Value: value}
	m.gauge[name] = &item
	log.Printf("SetGauge[%s]=%v\n", name, value)
}

func (m *MemStorage) IncCounter(name string, value int64) {
	if item := m.GetCounter(name); item != nil {
		item.Value += value
		log.Printf("IncCounter[%s]=%v\n", name, m.counter[name].Value)
		return
	}
	item := entity.CounterValue{Name: name, Value: value}
	m.counter[name] = &item
	log.Printf("IncCounter[%s]=%v\n", name, value)
}
