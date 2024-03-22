package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kosalnik/metrics/internal/models"
	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/storage"
)

func NewRestGetHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := io.ReadAll(req.Body)
		logrus.WithField("body", string(data)).Info("Handle Get")
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
		case "gauge":
			v, ok := s.GetGauge(m.ID)
			if !ok {
				http.NotFound(w, req)
				return
			}
			m.Value = &v
			if out, err := json.Marshal(m); err != nil {
				http.Error(w, `"internal error"`, http.StatusInternalServerError)
			} else {
				logrus.WithField("body", string(out)).Info("Handle Get Result")
				_, _ = w.Write(out)
			}
			return
		case "counter":
			v, ok := s.GetCounter(m.ID)
			if !ok {
				http.NotFound(w, req)
				return
			}
			m.Delta = &v
			if out, err := json.Marshal(m); err != nil {
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
		mType := chi.URLParam(req, "type")
		mName := chi.URLParam(req, "name")
		switch mType {
		case "gauge":
			v, ok := s.GetGauge(mName)
			if !ok {
				http.NotFound(w, req)
				return
			}
			res := fmt.Sprintf("%v", v)
			if _, err := w.Write([]byte(res)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		case "counter":
			v, ok := s.GetCounter(mName)
			if !ok {
				http.NotFound(w, req)
				return
			}
			res := fmt.Sprintf("%v", v)
			if _, err := w.Write([]byte(res)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
