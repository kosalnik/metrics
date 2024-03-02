package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kosalnik/metrics/internal/storage"
	"log"
	"net/http"
	"strconv"
)

type UpdateHandler struct {
	Storage storage.Storage
}

func NewUpdateHandler(s storage.Storage) *UpdateHandler {
	return &UpdateHandler{
		Storage: s,
	}
}

func (h *UpdateHandler) Handle(res http.ResponseWriter, req *http.Request) {
	mType := chi.URLParam(req, "type")
	mName := chi.URLParam(req, "name")
	mVal := chi.URLParam(req, "value")
	log.Println(mType, mName, mVal)
	switch mType {
	case "gauge":
		HandleUpdateGauge(h.Storage, mName, mVal)(res, req)
		return
	case "counter":
		HandleUpdateCounter(h.Storage, mName, mVal)(res, req)
		return
	}
	msg := fmt.Sprintf("bad request. wrong type %v", mType)
	http.Error(res, msg, http.StatusBadRequest)
}

func HandleUpdateGauge(s storage.Storage, name, value string) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if v, err := strconv.ParseFloat(value, 64); err != nil {
			http.Error(res, fmt.Sprintf("bad request. expected int64 [%s]", value), http.StatusBadRequest)
			return
		} else {
			s.SetGauge(name, v)
		}
	}
}

func HandleUpdateCounter(s storage.Storage, name, value string) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if v, err := strconv.ParseInt(value, 10, 64); err != nil {
			http.Error(res, fmt.Sprintf("bad request. expected float64 [%s]", value), http.StatusBadRequest)
			return
		} else {
			s.IncCounter(name, v)
		}
	}
}
