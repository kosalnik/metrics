package crypt

import (
	"bytes"
	"io"
	"net/http"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/infra/logger"
)

func HashCheckMiddleware(cfg config.Hash) func(next http.Handler) http.Handler {
	if cfg.Key == "" {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedHash := ExtractSign(r)
				logger.Logger.WithField("hash", expectedHash).Debug("Get Hash Header")

				if expectedHash != "" {
					defer r.Body.Close()
					b, err := io.ReadAll(r.Body)
					if err != nil {
						http.Error(w, "empty body", http.StatusInternalServerError)
						return
					}
					if !VerifySign(b, expectedHash, []byte(cfg.Key)) {
						http.Error(w, "verify hash fail", http.StatusBadRequest)
						return
					}
					req := r.Clone(r.Context())
					req.Body = io.NopCloser(bytes.NewReader(b))
					next.ServeHTTP(w, req)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
