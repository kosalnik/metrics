package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	c := http.Client{}
	return &SenderRest{
		client: &c,
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
	url := fmt.Sprintf("http://%s/update/", c.config.CollectorAddress)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		logrus.WithFields(logrus.Fields{"key": k, "val": v}).WithError(err).Errorf("send gauge. fail make request")
		return
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-Type", "application/json")
	r, err := c.client.Do(req)
	//logrus.WithFields(logrus.Fields{"url": url, "body": string(data)}).Info("send gauge.")
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
	vv := float64(v)
	m := models.Metrics{
		ID:    k,
		MType: "counter",
		Delta: &v,
		Value: &vv,
	}
	data, err := json.Marshal(m)
	if err != nil {
		logrus.WithFields(logrus.Fields{"key": k, "val": v}).WithError(err).Errorf("send counter. fail marshal")
		return
	}
	body := bytes.NewReader(data)
	url := fmt.Sprintf("http://%s/update/", c.config.CollectorAddress)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		logrus.WithFields(logrus.Fields{"key": k, "val": v}).WithError(err).Errorf("send gauge. fail make request")
		return
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-Type", "application/json")
	r, err := c.client.Do(req)
	logrus.WithFields(logrus.Fields{"url": url, "body": string(data)}).Info("send counter")
	if err != nil {
		logrus.Errorf("Fail push: %s", err.Error())
		return
	}
	response, err := io.ReadAll(r.Body)
	if err == nil {
		logrus.WithField("res", string(response)).Info("send counter result")
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			logrus.Errorf("fail close body. %s", err.Error())
		}
	}()
}
