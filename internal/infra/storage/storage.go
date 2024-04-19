package storage

import (
	"context"
	"time"

	"github.com/kosalnik/metrics/internal/models"
)

//go:generate mockgen -source=storage.go -destination=./mock/storage.go -package=mock
type Storage interface {
	GetGauge(ctx context.Context, name string) (*models.Metrics, error)
	SetGauge(ctx context.Context, name string, value float64) (*models.Metrics, error)
	GetCounter(ctx context.Context, name string) (*models.Metrics, error)
	IncCounter(ctx context.Context, name string, value int64) (*models.Metrics, error)
	UpsertAll(ctx context.Context, list []models.Metrics) error
	GetAll(ctx context.Context) ([]models.Metrics, error)
	Ping(ctx context.Context) error
	Close() error
	UpdatedAt() time.Time
}
