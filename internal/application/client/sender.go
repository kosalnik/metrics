package client

import (
	"context"

	"github.com/kosalnik/metrics/internal/models"
)

type Sender interface {
	SendGauge(k string, v float64)
	SendCounter(k string, v int64)
	SendBatch(ctx context.Context, list []models.Metrics) error
}
