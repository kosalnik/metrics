package client

import (
	"context"
	"io"
	"net/http"

	"github.com/kosalnik/metrics/internal/models"
)

//go:generate mockgen -source=sender.go -destination=./mock/senger.go -package=mock
type Sender interface {
	SendGauge(k string, v float64)
	SendCounter(k string, v int64)
	SendBatch(ctx context.Context, list []models.Metrics) error
}

type HttpSender interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}
