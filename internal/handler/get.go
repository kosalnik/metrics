package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kosalnik/metrics/internal/storage"
	"net/http"
)

type GetHandler struct {
	storage storage.Storage
}

func NewGetHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	h := GetHandler{
		storage: s,
	}
	return h.Handler()
}

func (h *GetHandler) Handler() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		mType := chi.URLParam(req, "type")
		mName := chi.URLParam(req, "name")
		switch mType {
		case "gauge":
			if !h.storage.HasGauge(mName) {
				http.NotFound(w, req)
				return
			}
			v := h.storage.GetGauge(mName)
			res := fmt.Sprintf("%v", v)
			_, err := w.Write([]byte(res))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		case "counter":
			if !h.storage.HasCounter(mName) {
				http.NotFound(w, req)
				return
			}
			v := h.storage.GetCounter(mName)
			res := fmt.Sprintf("%v", v)
			_, err := w.Write([]byte(res))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
