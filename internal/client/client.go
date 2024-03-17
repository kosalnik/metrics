package client

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/metric"
	"github.com/sirupsen/logrus"
)

type Client struct {
	mu        sync.Mutex
	client    *http.Client
	config    *config.Agent
	gauge     map[string]float64
	poolCount int64
}

func NewClient(config config.Agent) *Client {
	return &Client{
		client: http.DefaultClient,
		config: &config,
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
				c.sendGauge(k, v)
			}
			c.sendCounter("PoolCount", c.poolCount)
			c.sendGauge("RandomValue", rand.Float64())
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

func (c *Client) sendGauge(k string, v float64) {
	r, err := c.client.Post(fmt.Sprintf("http://%s/update/gauge/%s/%v", c.config.CollectorAddress, k, v), "text/plain", nil)
	if err != nil {
		logrus.Errorf("fail push. %s", err.Error())
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			logrus.Errorf("fail close body. %s", err.Error())
		}
	}()
}

func (c *Client) sendCounter(k string, v int64) {
	r, err := c.client.Post(fmt.Sprintf("http://%s/update/counter/%s/%v", c.config.CollectorAddress, k, v), "text/plain", nil)
	if err != nil {
		logrus.Errorf("Fail push: %s", err.Error())
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			logrus.Errorf("fail close body. %s", err.Error())
		}
	}()
}
