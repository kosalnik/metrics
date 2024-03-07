package client

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/metric"
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
	log.Printf("Pool interval: %d\n", c.config.PoolInterval)
	log.Printf("Report interval: %d\n", c.config.ReportInterval)
	log.Printf("Collector address: %s\n", c.config.CollectorAddress)
	go c.pool()
	c.push()
}

func (c *Client) push() {
	for {
		time.Sleep(time.Duration(c.config.ReportInterval) * time.Second)
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
		c.gauge = metric.GetMetrics()
		c.poolCount = c.poolCount + 1
		log.Println(c.poolCount)
		c.mu.Unlock()
		time.Sleep(time.Duration(c.config.PoolInterval) * time.Second)
	}
}

func (c *Client) sendGauge(k string, v float64) {
	r, err := c.client.Post(fmt.Sprintf("http://%s/update/gauge/%s/%v", c.config.CollectorAddress, k, v), "text/plain", nil)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("fail close body. %s", err.Error())
		}
	}()
	if err != nil {
		log.Printf("fail push. %s", err.Error())
	}
}

func (c *Client) sendCounter(k string, v int64) {
	r, err := c.client.Post(fmt.Sprintf("http://%s/update/counter/%s/%v", c.config.CollectorAddress, k, v), "text/plain", nil)
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("fail close body. %s", err.Error())
		}
	}()
	if err != nil {
		log.Printf("Fail push: %s", err.Error())
	}
}
