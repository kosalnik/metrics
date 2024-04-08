package storage

import (
	"context"

	"github.com/kosalnik/metrics/internal/models"
)

type Storage interface {
	GetGauge(ctx context.Context, name string) (float64, bool, error)
	SetGauge(ctx context.Context, name string, value float64) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, bool, error)
	IncCounter(ctx context.Context, name string, value int64) (int64, error)
	UpsertAll(ctx context.Context, list []models.Metrics) error
	GetAll(ctx context.Context) ([]models.Metrics, error)
	Ping(ctx context.Context) error
	Close() error
	Store(ctx context.Context, path string) error
	Recover(ctx context.Context, path string) error
}
