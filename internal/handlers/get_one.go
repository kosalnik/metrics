package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kosalnik/metrics/internal/infra/logger"
	"github.com/kosalnik/metrics/internal/infra/storage"
	"github.com/kosalnik/metrics/internal/models"
)

func NewRestGetHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := io.ReadAll(req.Body)
		logger.Logger.WithField("body", string(data)).Info("Handle Get")
		if err != nil {
			http.Error(w, `"Wrong data"`, http.StatusBadRequest)
			return
		}
		var m models.Metrics
		if err := json.Unmarshal(data, &m); err != nil {
			http.Error(w, `"Wrong json"`, http.StatusBadRequest)
			return
		}
		switch m.MType {
		case models.MGauge:
			v, err := s.GetGauge(req.Context(), m.ID)
			if err != nil {
				http.Error(w, `"fail get gauge"`, http.StatusInternalServerError)
				return
			}
			if v == nil {
				http.NotFound(w, req)
				return
			}
			if out, err := json.Marshal(v); err != nil {
				http.Error(w, `"internal error"`, http.StatusInternalServerError)
			} else {
				logger.Logger.WithField("body", string(out)).Info("Handle Get Result")
				_, _ = w.Write(out)
			}
			return
		case models.MCounter:
			v, err := s.GetCounter(req.Context(), m.ID)
			if err != nil {
				http.Error(w, `"fail get counter"`, http.StatusInternalServerError)
				return
			}
			if v == nil {
				http.NotFound(w, req)
				return
			}
			if out, err := json.Marshal(v); err != nil {
				http.Error(w, `"internal error"`, http.StatusInternalServerError)
			} else {
				_, _ = w.Write(out)
			}
			return
		}
		http.Error(w, `"not found"`, http.StatusNotFound)
	}
}

func NewGetHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		mType := models.MType(chi.URLParam(req, "type"))
		mName := chi.URLParam(req, "name")
		switch mType {
		case models.MGauge:
			v, err := s.GetGauge(req.Context(), mName)
			if err != nil {
				http.Error(w, `"fail get gauge"`, http.StatusInternalServerError)
				return
			}
			if v == nil {
				http.NotFound(w, req)
				return
			}
			res := fmt.Sprintf("%v", v.Value)
			if _, err := w.Write([]byte(res)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		case models.MCounter:
			v, err := s.GetCounter(req.Context(), mName)
			if err != nil {
				http.Error(w, `"fail get counter"`, http.StatusInternalServerError)
				return
			}
			if v == nil {
				http.NotFound(w, req)
				return
			}
			res := fmt.Sprintf("%v", v.Delta)
			if _, err := w.Write([]byte(res)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
