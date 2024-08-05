package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kosalnik/metrics/internal/crypt"
	"github.com/kosalnik/metrics/internal/log"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/models"
)

type SenderRest struct {
	client *HttpClientWrapper
	config *config.Agent
}

func NewSenderRest(config *config.Agent) Sender {
	mutators := []Mutator{
		SetRealIPMutator(),
		crypt.NewSignMutator([]byte(config.Hash.Key)),
	}
	if config.PublicKey != nil {
		mutators = append(mutators, crypt.NewCipherRequestMutator(crypt.NewEncoder(config.PublicKey, rand.Reader)))
	}
	return &SenderRest{
		client: NewHttpClient(WithMutators(mutators...)),
		config: config,
	}
}

func (c *SenderRest) SendGauge(k string, v float64) {
	m := models.Metrics{
		ID:    k,
		MType: models.MGauge,
		Value: v,
	}
	data, err := json.Marshal(m)
	if err != nil {
		log.Error().Str("key", k).Float64("val", v).Err(err).Msg("send gauge. fail marshal")
		return
	}
	body := bytes.NewReader(data)
	url := fmt.Sprintf("http://%s/update/", c.config.CollectorAddress)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Error().Str("key", k).Float64("val", v).Err(err).Msg("send gauge. fail make request")
		return
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-Type", "application/json")
	r, err := c.client.Do(req)
	if err != nil {
		log.Error().Str("key", k).Float64("val", v).Err(err).Msg("send gauge. fail post")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error().Err(err).Msg("fail close body.")
		}
	}()
}

func (c *SenderRest) SendCounter(k string, v int64) {
	vv := float64(v)
	m := models.Metrics{
		ID:    k,
		MType: models.MCounter,
		Delta: v,
		Value: vv,
	}
	data, err := json.Marshal(m)
	if err != nil {
		log.Error().Str("key", k).Int64("val", v).Err(err).Msg("send counter. fail marshal")
		return
	}
	body := bytes.NewReader(data)
	url := fmt.Sprintf("http://%s/update/", c.config.CollectorAddress)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Error().Str("key", k).Int64("val", v).Err(err).Msg("send gauge. fail make request")
		return
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-Type", "application/json")
	r, err := c.client.Do(req)
	log.Info().Str("url", url).Str("body", string(data)).Msg("send counter")
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

func (c *SenderRest) SendBatch(ctx context.Context, list []models.Metrics) error {
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
