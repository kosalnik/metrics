package client

import (
	"math/rand"
	"sync"
	"time"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/infra/metric"
	"github.com/sirupsen/logrus"
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

func (c *Client) Run() {
	logrus.Infof("Poll interval: %d", c.config.PollInterval)
	logrus.Infof("Report interval: %d", c.config.ReportInterval)
	logrus.Infof("Collector address: %s", c.config.CollectorAddress)
	go c.poll()
	c.push()
}

func (c *Client) push() {
	for {
		time.Sleep(time.Duration(c.config.ReportInterval) * time.Second)
		c.mu.Lock()
		logrus.Info("Push")
		if c.gauge != nil {
			for k, v := range c.gauge {
				c.sender.SendGauge(k, v)
			}
		}
		c.sender.SendCounter("PollCount", c.pollCount)
		c.sender.SendGauge("RandomValue", rand.Float64())
		c.mu.Unlock()
	}
}

func (c *Client) poll() {
	for {
		c.mu.Lock()
		logrus.Debug("Poll")
		c.gauge = metric.GetMetrics()
		c.pollCount = c.pollCount + 1
		logrus.Debugf("PollCount=%d", c.pollCount)
		c.mu.Unlock()
		time.Sleep(time.Duration(c.config.PollInterval) * time.Second)
	}
}
