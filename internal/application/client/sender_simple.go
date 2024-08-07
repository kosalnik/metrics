package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kosalnik/metrics/internal/log"
	"github.com/kosalnik/metrics/internal/models"

	"github.com/kosalnik/metrics/internal/config"
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
		log.Error().Err(err).Msg("fail push")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error().Err(err).Msg("fail close body")
		}
	}()
}

func (c *SenderSimple) SendCounter(k string, v int64) {
	r, err := c.client.Post(fmt.Sprintf("http://%s/update/counter/%s/%v", c.config.CollectorAddress, k, v), "text/plain", nil)
	if err != nil {
		log.Error().Err(err).Msg("Fail push")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error().Err(err).Msg("fail close body")
		}
	}()
}

func (c *SenderSimple) SendBatch(ctx context.Context, list []models.Metrics) error {
	if len(list) == 0 {
		return nil
	}
	data, err := json.Marshal(list)
	if err != nil {
		log.Error().Any("list", list).Err(err).Msg("fail send batch. fail marshal")

		return err
	}
	body := bytes.NewReader(data)
	url := fmt.Sprintf("http://%s/updates/", c.config.CollectorAddress)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		log.Error().Err(err).Msg("send batch. fail make request")

		return err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-Type", "application/json")
	r, err := c.client.Do(req)
	log.Info().Str("url", url).Str("body", string(data)).Msg("send counter")
	if err != nil {
		log.Error().Err(err).Msg("Fail push")

		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error().Err(err).Msg("fail close body")
		}
	}()

	return nil
}
