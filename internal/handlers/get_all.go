package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/kosalnik/metrics/internal/storage"
)

func NewGetAllHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		items := s.GetPlain()
		var res []string
		for k, v := range items {
			res = append(res, fmt.Sprintf("%s = %s", k, v))
		}
		sort.Strings(res)
		if _, err := fmt.Fprint(w, strings.Join(res, "\n")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
