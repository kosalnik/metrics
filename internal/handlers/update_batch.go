package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/kosalnik/metrics/internal/log"
	"github.com/kosalnik/metrics/internal/models"
	"github.com/kosalnik/metrics/internal/storage"
)

func NewUpdateBatchHandler(s storage.BatchInserter) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		data, err := io.ReadAll(req.Body)
		log.Debug().Str("data", string(data)).Msg("Handle batch update")
		if err != nil {
			http.Error(res, `"Wrong data"`, http.StatusBadRequest)
			return
		}
		var mList []models.Metrics
		if err := json.Unmarshal(data, &mList); err != nil {
			http.Error(res, `"Wrong json"`, http.StatusBadRequest)
			return
		}
		if err := s.UpsertAll(req.Context(), mList); err != nil {
			http.Error(res, `"fail upsert"`, http.StatusInternalServerError)
			return
		}
	}
}
