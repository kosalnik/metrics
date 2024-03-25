package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/kosalnik/metrics/internal/infra/storage"
	"github.com/sirupsen/logrus"
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
				t = append(t, v.String())
			}
			sort.Strings(t)
			data = []byte(strings.Join(t, "\n"))
		}
		if _, err := w.Write(data); err != nil {
			logrus.WithError(err).Error("fail write response")
		}
	}
}
