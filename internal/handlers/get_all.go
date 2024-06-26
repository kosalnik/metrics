package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/kosalnik/metrics/internal/infra/logger"
	"github.com/kosalnik/metrics/internal/infra/storage"
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
			logger.Logger.WithError(err).Error("fail get all")
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
				logger.Logger.WithField("v", v).Info("asdf")
				t[i] = v.String()
			}
			sort.Strings(t)
			data = []byte(strings.Join(t, "\n"))
		}
		if _, err := w.Write(data); err != nil {
			logger.Logger.WithError(err).Error("fail write response")
		}
	}
}
