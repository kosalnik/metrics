package handler

import (
	"fmt"
	"log"
	"metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

func HandleUpdateGauge(s storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "bad request method", http.StatusBadRequest)
			return

		}
		p, err := parsePath(req.URL.Path)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if v, err := strconv.ParseFloat(p[1], 64); err != nil {
			http.Error(res, fmt.Sprintf("bad request. expected int64 [%s]", p[1]), http.StatusBadRequest)
			return
		} else {
			s.SetGauge(p[0], v)
		}

		log.Printf("Handled updateGauge: %v", p)
	}
}

func HandleUpdateCounter(s storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		p, err := parsePath(req.URL.Path)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if v, err := strconv.ParseInt(p[1], 10, 64); err != nil {
			http.Error(res, fmt.Sprintf("bad request. expected float64 [%s]", p[1]), http.StatusBadRequest)
			return
		} else {
			s.IncCounter(p[0], v)
		}
		log.Printf("Handled updateCounter: %v", p)
	}
}

func parsePath(u string) ([]string, error) {
	p := strings.Split(u, "/")
	if len(p) != 5 {
		return nil, fmt.Errorf("Wrong request %v %v", len(p), p)
	}
	return p[3:], nil
}
