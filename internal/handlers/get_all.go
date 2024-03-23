package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kosalnik/metrics/internal/storage"
)

func NewGetAllHandler(s storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		accept := req.Header.Get("Accept")
		isJSON := accept != "" && strings.Contains(accept, "application/json")
		if isJSON {
			w.Header().Set("Content-Type", "application/json")
		} else {
			w.Header().Set("Content-Type", "text/html")
		}
		items := s.GetAll()
		var err error
		var data []byte
		if isJSON {
			data, err = json.Marshal(items)
			if err != nil {
				http.Error(w, `"fail marshal"`, http.StatusInternalServerError)
			}
		} else {
			var t []string
			for _, v := range items {
				if v.MType == "counter" {
					t = append(t, fmt.Sprintf("%s = %v", v.ID, *v.Delta))
				} else {
					t = append(t, fmt.Sprintf("%s = %v", v.ID, *v.Value))
				}
			}
			data = []byte(strings.Join(t, "\n"))
		}
		w.Write(data)
	}
}
