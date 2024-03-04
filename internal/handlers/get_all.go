package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/kosalnik/metrics/internal/storage"
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
		var res = []string{}
		for k, v := range items {
			res = append(res, fmt.Sprintf("%s = %s", k, v))
		}
		sort.Strings(res)
		_, err := fmt.Fprint(w, strings.Join(res, "\n"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
