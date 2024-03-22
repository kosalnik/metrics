package client

import (
	"fmt"
	"net/http"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/sirupsen/logrus"
)

type SenderSimple struct {
	client *http.Client
	config *config.Agent
}

func NewSenderSimple(config *config.Agent) Sender {
	return &SenderSimple{
		client: http.DefaultClient,
		config: config,
	}
}

func (c *SenderSimple) SendGauge(k string, v float64) {
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

func (c *SenderSimple) SendCounter(k string, v int64) {
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
