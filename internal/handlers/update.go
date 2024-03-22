package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kosalnik/metrics/internal/models"
	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/storage"
)

func NewRestUpdateHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, "Wrong data", http.StatusBadRequest)
			return
		}
		var m models.Metrics
		if err := json.Unmarshal(data, &m); err != nil {
			http.Error(res, "Wrong json", http.StatusBadRequest)
			return
		}
		logrus.Debugf("Handle %s[%s]=(%f|%d)", m.MType, m.ID, *m.Value, *m.Delta)
		switch m.MType {
		case "gauge":
			r := s.SetGauge(m.ID, *m.Value)
			m.Value = &r
			if out, err := json.Marshal(m); err != nil {
				http.Error(res, "internal error", http.StatusInternalServerError)
			} else {
				_, _ = res.Write(out)
			}
			return
		case "counter":
			r := s.IncCounter(m.ID, *m.Delta)
			m.Delta = &r
			if out, err := json.Marshal(m); err != nil {
				http.Error(res, "internal error", http.StatusInternalServerError)
			} else {
				_, _ = res.Write(out)
			}
			return
		}
		http.Error(res, "bad request", http.StatusBadRequest)
	}
}

func NewUpdateHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		mType := chi.URLParam(req, "type")
		mName := chi.URLParam(req, "name")
		mVal := chi.URLParam(req, "value")
		logrus.Debugf("Handle %s[%s]=%s", mType, mName, mVal)
		switch mType {
		case "gauge":
			v, err := strconv.ParseFloat(mVal, 64)
			if err != nil {
				http.Error(res, "bad request", http.StatusBadRequest)
				return
			}
			s.SetGauge(mName, v)
			return
		case "counter":
			v, err := strconv.ParseInt(mVal, 10, 64)
			if err != nil {
				http.Error(res, "bad request", http.StatusBadRequest)
				return
			}
			s.IncCounter(mName, v)
			return
		}
		http.Error(res, "bad request", http.StatusBadRequest)
	}
}
