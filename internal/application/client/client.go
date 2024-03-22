package client

import (
	"math/rand"
	"sync"
	"time"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/metric"
	"github.com/sirupsen/logrus"
)

type Client struct {
	mu        sync.Mutex
	sender    Sender
	config    *config.Agent
	gauge     map[string]float64
	poolCount int64
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
	logrus.Infof("Pool interval: %d", c.config.PoolInterval)
	logrus.Infof("Report interval: %d", c.config.ReportInterval)
	logrus.Infof("Collector address: %s", c.config.CollectorAddress)
	go c.pool()
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
			c.sender.SendCounter("PoolCount", c.poolCount)
			c.sender.SendGauge("RandomValue", rand.Float64())
		}
		c.mu.Unlock()
	}
}

func (c *Client) pool() {
	for {
		c.mu.Lock()
		logrus.Debug("Pool")
		c.gauge = metric.GetMetrics()
		c.poolCount = c.poolCount + 1
		logrus.Debugf("PoolCount=%d", c.poolCount)
		c.mu.Unlock()
		time.Sleep(time.Duration(c.config.PoolInterval) * time.Second)
	}
}
