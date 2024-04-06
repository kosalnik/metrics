package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/infra/storage"
	"github.com/kosalnik/metrics/internal/models"
)

func NewRestUpdateHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		data, err := io.ReadAll(req.Body)
		logrus.Debugf("Handle %s", data)
		if err != nil {
			http.Error(res, `"Wrong data"`, http.StatusBadRequest)
			return
		}
		var m models.Metrics
		if err := json.Unmarshal(data, &m); err != nil {
			http.Error(res, `"Wrong json"`, http.StatusBadRequest)
			return
		}
		switch m.MType {
		case models.MGauge:
			r := s.SetGauge(m.ID, *m.Value)
			m.Value = &r
			if out, err := json.Marshal(m); err != nil {
				http.Error(res, `"internal error"`, http.StatusInternalServerError)
			} else {
				_, _ = res.Write(out)
			}
			return
		case models.MCounter:
			r := s.IncCounter(m.ID, *m.Delta)
			m.Delta = &r
			if out, err := json.Marshal(m); err != nil {
				http.Error(res, `"internal error"`, http.StatusInternalServerError)
			} else {
				_, _ = res.Write(out)
			}
			return
		}
		http.Error(res, `"not found"`, http.StatusNotFound)
	}
}

func NewUpdateHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain")
		mType := models.MType(chi.URLParam(req, "type"))
		mName := chi.URLParam(req, "name")
		mVal := chi.URLParam(req, "value")
		logrus.Debugf("Handle %s[%s]=%s", mType, mName, mVal)
		switch mType {
		case models.MGauge:
			v, err := strconv.ParseFloat(mVal, 64)
			if err != nil {
				http.Error(res, "bad request", http.StatusBadRequest)

				return
			}
			r := s.SetGauge(mName, v)
			if _, err := res.Write([]byte(fmt.Sprintf("%f", r))); err != nil {
				logrus.WithError(err).Error("fail write response")
			}

			return
		case models.MCounter:
			v, err := strconv.ParseInt(mVal, 10, 64)
			if err != nil {
				http.Error(res, "bad request", http.StatusBadRequest)
				return
			}
			r := s.IncCounter(mName, v)
			if _, err := res.Write([]byte(fmt.Sprintf("%d", r))); err != nil {
				logrus.WithError(err).Error("fail write response")
			}

			return
		}
		http.Error(res, "bad request", http.StatusBadRequest)
	}
}
