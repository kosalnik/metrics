// Package client содержит реализацию клиента к коллектору метрик.
// Создаётся с использованием NewClient(). Запускается методом Run().
// При старте запускает два параллельных цикла. Один собирает метрики раз в PoolInterval секунд.
// Второй цикл отсылает собранные в последний раз метрики коллектору раз в ReportInterval секунд.
package client

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/log"
	"github.com/kosalnik/metrics/internal/metric"
	"github.com/kosalnik/metrics/internal/models"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	sender    Sender
	config    *config.Agent
	gauge     map[string]float64
	pollCount int64
	mu        sync.Mutex
}

func NewClient(ctx context.Context, config *config.Agent) *Client {
	return &Client{
		config: config,
		sender: NewSenderPool(
			ctx,
			NewSenderRest(config),
			int(config.RateLimit),
		),
	}
}

func (c *Client) Run(ctx context.Context) error {
	g := errgroup.Group{}
	g.Go(func() error {
		return c.RunMetricsSender(ctx)
	})
	g.Go(func() error {
		return c.RunMetricsPoller(ctx)
	})
	return g.Wait()
}

// RunMetricsSender Запустить процесс периодической отправки метрик в коллектор
func (c *Client) RunMetricsSender(ctx context.Context) error {
	log.Info().
		Int64("Report interval", c.config.ReportInterval).
		Str("Collector address", c.config.CollectorAddress).
		Msg("Running agent")
	return c.push(ctx)
}

// RunMetricsPoller Запустить процесс периодического сбора метрик
func (c *Client) RunMetricsPoller(ctx context.Context) error {
	log.Info().
		Int64("Poll interval", c.config.PollInterval).
		Msg("Running metrics poll process")
	return c.poll(ctx)
}

func (c *Client) Shutdown(ctx context.Context) {
	log.Info().Msg("Shutdown service Client")
	if err := c.sender.SendBatch(ctx, c.collectMetrics()); err != nil {
		log.Error().Err(err).Msg("fail push")
	}
	log.Info().Msg("Shutdown service Client completed")
}

func (c *Client) push(ctx context.Context) error {
	tick := time.NewTicker(time.Duration(c.config.ReportInterval) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			log.Info().Msg("Push")
			if err := c.sender.SendBatch(ctx, c.collectMetrics()); err != nil {
				log.Error().Err(err).Msg("fail push")
			}
		}
	}
}

func (c *Client) collectMetrics() []models.Metrics {
	c.mu.Lock()
	defer c.mu.Unlock()

	list := make([]models.Metrics, len(c.gauge)+2)
	i := 0
	if c.gauge != nil {
		for k, v := range c.gauge {
			kk := k
			vv := v
			list[i] = models.Metrics{ID: kk, MType: models.MGauge, Value: vv}
			i++
		}
	}

	vv := c.pollCount
	list[i] = models.Metrics{ID: "PollCount", MType: models.MCounter, Delta: vv}

	rv := rand.Float64()
	list[i+1] = models.Metrics{ID: "RandomValue", MType: models.MGauge, Value: rv}

	return list
}

func (c *Client) poll(ctx context.Context) error {
	tick := time.NewTicker(time.Duration(c.config.PollInterval) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			if err := c.pollMetrics(ctx); err != nil {
				log.Error().Err(err).Msg("poll error")
			}
		}
	}
}

func (c *Client) pollMetrics(ctx context.Context) error {
	var err error
	c.mu.Lock()
	defer c.mu.Unlock()
	c.gauge, err = metric.GetMetrics(ctx)
	if err != nil {
		return err
	}
	c.pollCount = c.pollCount + 1
	log.Debug().Int64("count", c.pollCount).Msg("PollCount")

	return nil
}
