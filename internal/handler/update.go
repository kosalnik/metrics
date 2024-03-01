package handler

import (
	"fmt"
	"github.com/kosalnik/metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

type UpdateHandler struct {
	Storage storage.Storage
}

func NewUpdateHandler(s storage.Storage) *UpdateHandler {
	return &UpdateHandler{
		Storage: s,
	}
}

func (h *UpdateHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "bad request method", http.StatusBadRequest)
		return
	}
	p := strings.Split(req.URL.Path, "/")
	if len(p) != 5 {
		msg := fmt.Sprintf("bad request %v %v", len(p), p)
		http.Error(res, msg, http.StatusNotFound)
		return
	}
	mType, mName, mVal := p[2], p[3], p[4]
	switch mType {
	case "gauge":
		HandleUpdateGauge(h.Storage, mName, mVal)(res, req)
		return
	case "counter":
		HandleUpdateCounter(h.Storage, mName, mVal)(res, req)
		return
	}
	msg := fmt.Sprintf("bad request. wrong type %v", mType)
	http.Error(res, msg, http.StatusNotFound)
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
