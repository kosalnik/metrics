package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kosalnik/metrics/internal/storage"
)

func NewRestGetHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

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
