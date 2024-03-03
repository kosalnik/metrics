package handlers

import (
	"fmt"
	"github.com/kosalnik/metrics/internal/storage"
	"net/http"
	"sort"
	"strings"
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
		fmt.Println(items)
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
