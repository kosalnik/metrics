package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

func NewPingHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			http.Error(w, "", http.StatusInternalServerError)

			return
		}
	}
}
