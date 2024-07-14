// Package storage содержит интерфейсы, которые должен реализовывать storage.
package storage

import (
	"context"
	"io"
	"time"

	"github.com/kosalnik/metrics/internal/models"
)

//go:generate mockgen -source=storage.go -destination=./mock/storage.go -package=mock
type Storage interface {
	Dumper
	BatchInserter
	UpdateAwarer
	Pinger
	GetGauge(ctx context.Context, name string) (*models.Metrics, error)
	SetGauge(ctx context.Context, name string, value float64) (*models.Metrics, error)
	GetCounter(ctx context.Context, name string) (*models.Metrics, error)
	IncCounter(ctx context.Context, name string, value int64) (*models.Metrics, error)
	io.Closer
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Dumper interface {
	GetAll(ctx context.Context) ([]models.Metrics, error)
}

type BatchInserter interface {
	UpsertAll(ctx context.Context, list []models.Metrics) error
}

type UpdateAwarer interface {
	UpdatedAt() time.Time
}
