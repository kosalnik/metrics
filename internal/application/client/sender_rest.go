package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/models"
	"github.com/sirupsen/logrus"
)

type SenderRest struct {
	client *http.Client
	config *config.Agent
}

func NewSenderRest(config *config.Agent) Sender {
	return &SenderSimple{
		client: http.DefaultClient,
		config: config,
	}
}

func (c *SenderRest) SendGauge(k string, v float64) {
	m := models.Metrics{
		ID:    k,
		MType: "gauge",
		Value: &v,
	}
	data, err := json.Marshal(m)
	if err != nil {
		logrus.WithFields(logrus.Fields{"key": k, "val": v}).WithError(err).Errorf("send gauge. fail marshal")
		return
	}
	body := bytes.NewReader(data)
	r, err := c.client.Post(fmt.Sprintf("http://%s/update/", c.config.CollectorAddress), "application/json", body)
	if err != nil {
		logrus.WithFields(logrus.Fields{"key": k, "val": v}).WithError(err).Errorf("send gauge. fail post")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			logrus.Errorf("fail close body. %s", err.Error())
		}
	}()
}

func (c *SenderRest) SendCounter(k string, v int64) {
	m := models.Metrics{
		ID:    k,
		MType: "counter",
		Delta: &v,
	}
	data, err := json.Marshal(m)
	if err != nil {
		logrus.WithFields(logrus.Fields{"key": k, "val": v}).WithError(err).Errorf("send gauge. fail marshal")
		return
	}
	body := bytes.NewReader(data)
	r, err := c.client.Post(fmt.Sprintf("http://%s/update/", c.config.CollectorAddress), "application/json", body)
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
