package client

import (
	"fmt"
	"github.com/kosalnik/metrics/internal/config"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	mu        sync.Mutex
	client    *http.Client
	config    *config.ClientConfig
	gauge     map[string]float64
	poolCount int64
}

func NewClient(config config.ClientConfig) *Client {
	return &Client{
		client: http.DefaultClient,
		config: &config,
	}
}

func (c *Client) Run() {
	go c.pool()
	c.push()
}

func (c *Client) push() {
	for {
		time.Sleep(c.config.ReportInterval)
		c.mu.Lock()
		log.Println("Push")
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
		log.Println("Pool")
		c.gauge = getMetrics()
		c.poolCount = c.poolCount + 1
		log.Println(c.poolCount)
		c.mu.Unlock()
		time.Sleep(c.config.PoolInterval)
	}
}

func (c *Client) sendGauge(k string, v float64) {
	r, err := c.client.Post(fmt.Sprintf("%s/update/gauge/%s/%v", c.config.CollectorAddress, k, v), "text/plain", nil)
	if err != nil {
		log.Printf("fail push. %s", err.Error())
	}
	if err := r.Body.Close(); err != nil {
		log.Printf("fail close body. %s", err.Error())
	}
}

func (c *Client) sendCounter(k string, v int64) {
	r, err := c.client.Post(fmt.Sprintf("%s/update/counter/%s/%v", c.config.CollectorAddress, k, v), "text/plain", nil)
	if err != nil {
		log.Printf("Fail push: %s", err.Error())
	}
	if err := r.Body.Close(); err != nil {
		log.Printf("fail close body. %s", err.Error())
	}
}
