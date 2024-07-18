package crypt

import (
	"bytes"
	"io"
	"net/http"
)

//go:generate mockgen -source=cipher_middleware.go -destination=./mock/decoder.go -package=mock
type Decoder interface {
	Decode([]byte) ([]byte, error)
}

func CipherMiddleware(decoder Decoder) func(next http.Handler) http.Handler {
	if decoder == nil {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				ciphertext, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "empty body", http.StatusInternalServerError)
					return
				}
				plaintext, err := decoder.Decode(ciphertext)
				if err != nil {
					http.Error(w, "wrong request", http.StatusBadRequest)
					return
				}
				req := r.Clone(r.Context())
				req.Body = io.NopCloser(bytes.NewReader(plaintext))
				next.ServeHTTP(w, req)
			},
		)
	}
}
