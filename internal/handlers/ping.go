package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/kosalnik/metrics/internal/infra/storage"
)

func NewPingHandler(db storage.Storage) func(res http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			http.Error(w, "", http.StatusInternalServerError)

			return
		}
	}
}
