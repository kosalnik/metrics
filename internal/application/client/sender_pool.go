package client

import (
	"context"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/infra/logger"
	"github.com/kosalnik/metrics/internal/models"
)

type SenderPool struct {
	client Sender
	jobs   chan func()
}

var _ Sender = &SenderPool{}

func NewSenderPool(ctx context.Context, cfg *config.Agent) *SenderPool {
	p := &SenderPool{
		client: NewSenderRest(cfg),
		jobs:   make(chan func()),
	}

	for w := 1; w <= int(cfg.RateLimit); w++ {
		logger.Logger.Infof("Start worker %d", w)
		go p.worker(p.jobs)
	}

	go func() {
		<-ctx.Done()
		close(p.jobs)
	}()

	return p
}

func (p *SenderPool) SendGauge(k string, v float64) {
	p.jobs <- func() {
		p.client.SendGauge(k, v)
	}
}

func (p *SenderPool) SendCounter(k string, v int64) {
	p.jobs <- func() {
		p.client.SendCounter(k, v)
	}
}

func (p *SenderPool) SendBatch(ctx context.Context, list []models.Metrics) error {
	logger.Logger.Debug("Push job")
	p.jobs <- func() {
		if err := p.client.SendBatch(ctx, list); err != nil {
			logger.Logger.WithError(err).Error("Send batch failed")
		}
	}

	return nil
}

func (*SenderPool) worker(jobs <-chan func()) {
	for f := range jobs {
		logger.Logger.Debug("worker receive job")
		f()
	}
	logger.Logger.Info("Stop worker")
}
