package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/kosalnik/metrics/internal/infra/logger"
	"github.com/kosalnik/metrics/internal/infra/storage"
	"github.com/kosalnik/metrics/internal/models"
)

func NewRestUpdateHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		data, err := io.ReadAll(req.Body)
		logger.Logger.Debugf("Handle %s", data)
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
			r, err := s.SetGauge(req.Context(), m.ID, m.Value)
			if err != nil {
				http.Error(res, `"fail set gauge"`, http.StatusInternalServerError)
				return
			}
			if out, err := json.Marshal(r); err != nil {
				http.Error(res, `"internal error"`, http.StatusInternalServerError)
			} else {
				_, _ = res.Write(out)
			}
			return
		case models.MCounter:
			r, err := s.IncCounter(req.Context(), m.ID, m.Delta)
			if err != nil {
				http.Error(res, `"fail inc counter"`, http.StatusInternalServerError)
				return
			}
			if out, err := json.Marshal(r); err != nil {
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
		logger.Logger.Debugf("Handle %s[%s]=%s", mType, mName, mVal)
		switch mType {
		case models.MGauge:
			v, err := strconv.ParseFloat(mVal, 64)
			if err != nil {
				http.Error(res, "bad request", http.StatusBadRequest)

				return
			}
			r, err := s.SetGauge(req.Context(), mName, v)
			if err != nil {
				http.Error(res, `"fail set gauge"`, http.StatusInternalServerError)
				return
			}
			if _, err := res.Write([]byte(fmt.Sprintf("%f", r.Value))); err != nil {
				logger.Logger.WithError(err).Error("fail write response")
			}

			return
		case models.MCounter:
			v, err := strconv.ParseInt(mVal, 10, 64)
			if err != nil {
				http.Error(res, "bad request", http.StatusBadRequest)
				return
			}
			r, err := s.IncCounter(req.Context(), mName, v)
			if err != nil {
				logger.Logger.WithError(err).Error("fail inc counter")
				http.Error(res, `"fail inc counter"`, http.StatusInternalServerError)
				return
			}
			if _, err := res.Write([]byte(fmt.Sprintf("%d", r.Delta))); err != nil {
				logger.Logger.WithError(err).Error("fail write response")
			}

			return
		}
		http.Error(res, "bad request", http.StatusBadRequest)
	}
}
