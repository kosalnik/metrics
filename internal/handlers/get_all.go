// Package handlers contains handlers.
// Методы этого пакета создают хендлеры для http.Server
package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/kosalnik/metrics/internal/log"
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
		items, err := s.GetAll(req.Context())
		if err != nil {
			log.Error().Err(err).Msg("fail get all")
			http.Error(w, `"fail get all"`, http.StatusInternalServerError)

			return
		}
		if len(items) == 0 {
			w.Write([]byte(`[]`))
			return
		}
		var data []byte
		if isJSON {
			data, err = json.Marshal(items)
			if err != nil {
				http.Error(w, `"fail marshal"`, http.StatusInternalServerError)

				return
			}
		} else {
			t := make([]string, len(items))
			for i, v := range items {
				log.Info().Any("v", v).Msg("get all metrics")
				t[i] = v.String()
			}
			sort.Strings(t)
			data = []byte(strings.Join(t, "\n"))
		}
		if _, err := w.Write(data); err != nil {
			log.Error().Err(err).Msg("fail write response")
		}
	}
}
