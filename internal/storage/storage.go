package storage

import "github.com/kosalnik/metrics/internal/models"

type Storage interface {
	GetGauge(name string) (float64, bool)
	SetGauge(name string, value float64) float64
	GetCounter(name string) (int64, bool)
	IncCounter(name string, value int64) int64
	GetAll() []models.Metrics
}
