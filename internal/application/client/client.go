package client

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/infra/logger"
	"github.com/kosalnik/metrics/internal/infra/metric"
	"github.com/kosalnik/metrics/internal/models"
)

type Client struct {
	mu        sync.Mutex
	sender    Sender
	config    *config.Agent
	gauge     map[string]float64
	pollCount int64
}

func NewClient(config config.Agent) *Client {
	return &Client{
		config: &config,
		sender: NewSenderRest(
			&config,
		),
	}
}

func (c *Client) Run(ctx context.Context) {
	logger.Logger.Infof("Poll interval: %d", c.config.PollInterval)
	logger.Logger.Infof("Report interval: %d", c.config.ReportInterval)
	logger.Logger.Infof("Collector address: %s", c.config.CollectorAddress)
	go c.poll(ctx)
	c.push(ctx)
}

func (c *Client) push(ctx context.Context) {
	tick := time.NewTicker(time.Duration(c.config.ReportInterval) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			logger.Logger.Info("Push")
			if err := c.sender.SendBatch(ctx, c.collectMetrics()); err != nil {
				logger.Logger.WithError(err).Error("fail push")
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
			list[i] = models.Metrics{ID: kk, MType: models.MGauge, Value: &vv}
			i++
		}
	}

	vv := c.pollCount
	list[i] = models.Metrics{ID: "PollCount", MType: models.MCounter, Delta: &vv}

	rv := rand.Float64()
	list[i+1] = models.Metrics{ID: "RandomValue", MType: models.MGauge, Value: &rv}

	return list
}

func (c *Client) poll(ctx context.Context) {
	tick := time.NewTicker(time.Duration(c.config.PollInterval) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			c.pollMetrics()
		}
	}
}

func (c *Client) pollMetrics() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.gauge = metric.GetMetrics()
	c.pollCount = c.pollCount + 1
	logger.Logger.WithField("count", c.pollCount).Debug("PollCount")
}
