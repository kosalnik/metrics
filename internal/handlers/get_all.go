package handlers

import (
	"fmt"
	"github.com/kosalnik/metrics/internal/storage"
	"net/http"
)

type GetAllHandler struct {
	storage storage.Storage
}

func NewGetAllHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	h := GetAllHandler{
		storage: s,
	}
	return h.Handler()
}

func (h *GetAllHandler) Handler() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		items := h.storage.GetPlain()
		for k, v := range items {
			_, err := fmt.Fprintf(w, "%s = %s\n", k, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
