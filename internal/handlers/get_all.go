package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kosalnik/metrics/internal/storage"
)

func NewGetAllHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		accept := req.Header.Get("Accept")
		isJSON := strings.Contains(accept, "application/json")
		if isJSON {
			w.Header().Set("Content-Type", "application/json")
		} else {
			w.Header().Set("Content-Type", "text/html")
		}
		items := s.GetAll()
		data, err := json.Marshal(items)
		if err != nil {
			http.Error(w, `"fail marshal"`, http.StatusInternalServerError)
		}
		w.Write(data)
	}
}
