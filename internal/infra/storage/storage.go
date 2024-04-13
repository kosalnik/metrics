package storage

import (
	"context"
	"time"

	"github.com/kosalnik/metrics/internal/models"
)

//go:generate mockgen -source=storage.go -destination=./mock/storage.go -package=mock
type Storage interface {
	GetGauge(ctx context.Context, name string) (float64, bool, error)
	SetGauge(ctx context.Context, name string, value float64) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, bool, error)
	IncCounter(ctx context.Context, name string, value int64) (int64, error)
	UpsertAll(ctx context.Context, list []models.Metrics) error
	GetAll(ctx context.Context) ([]models.Metrics, error)
	Ping(ctx context.Context) error
	Close() error
	UpdatedAt() time.Time
}
