package client

import (
	"context"

	"github.com/kosalnik/metrics/internal/log"
	"github.com/kosalnik/metrics/internal/models"
)

type SenderPool struct {
	client Sender
	jobs   chan func()
}

var _ Sender = &SenderPool{}

func NewSenderPool(ctx context.Context, client Sender, num int) *SenderPool {
	p := &SenderPool{
		client: client,
		jobs:   make(chan func()),
	}

	for w := 1; w <= int(num); w++ {
		log.Info().Int("id", w).Msg("Start worker")
		go p.worker(p.jobs)
	}

	go func() {
		<-ctx.Done()
		close(p.jobs)
	}()

	return p
}

func (p *SenderPool) Shutdown(_ context.Context) {
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
	log.Debug().Msg("Push job")
	p.jobs <- func() {
		if err := p.client.SendBatch(ctx, list); err != nil {
			log.Error().Err(err).Msg("Send batch failed")
		}
	}

	return nil
}

func (*SenderPool) worker(jobs <-chan func()) {
	for f := range jobs {
		log.Debug().Msg("worker receive job")
		f()
	}
	log.Info().Msg("Stop worker")
}
