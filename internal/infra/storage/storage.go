package storage

import (
	"context"
	"time"

	"github.com/kosalnik/metrics/internal/models"
)

//go:generate mockgen -source=storage.go -destination=./mock/storage.go -package=mock
type Storage interface {
	Dumper
	Recoverer
	UpdateAwarer
	GetGauge(ctx context.Context, name string) (*models.Metrics, error)
	SetGauge(ctx context.Context, name string, value float64) (*models.Metrics, error)
	GetCounter(ctx context.Context, name string) (*models.Metrics, error)
	IncCounter(ctx context.Context, name string, value int64) (*models.Metrics, error)
	Ping(ctx context.Context) error
	Close() error
}

type Dumper interface {
	GetAll(ctx context.Context) ([]models.Metrics, error)
}

type Recoverer interface {
	UpsertAll(ctx context.Context, list []models.Metrics) error
}

type UpdateAwarer interface {
	UpdatedAt() time.Time
}
