package crypt_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/infra/crypt"
	"github.com/stretchr/testify/require"
)

func TestHashCheckMiddleware(t *testing.T) {
	tests := map[string]struct {
		cfg    config.Hash
		body   string
		header string
		want   int
	}{
		"Success. Empty Key": {
			cfg:    config.Hash{Key: ""},
			body:   `asdf`,
			header: `f553f2c73c8`,
			want:   http.StatusOK,
		},
		"Success. With Key": {
			cfg:    config.Hash{Key: "Secret"},
			body:   `Hello`,
			header: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			want:   http.StatusOK,
		},
		"Invalid hash. Other Key": {
			cfg:    config.Hash{Key: "asd"},
			body:   `Hello`,
			header: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			want:   http.StatusBadRequest,
		},
		"Invalid hash. Other body": {
			cfg:    config.Hash{Key: "Secret"},
			body:   `asd`,
			header: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			want:   http.StatusBadRequest,
		},
		"Invalid hash. Other hash": {
			cfg:    config.Hash{Key: "Secret"},
			body:   `Hello`,
			header: `asdf`,
			want:   http.StatusBadRequest,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mw := crypt.HashCheckMiddleware(tt.cfg)
			mux := http.NewServeMux()
			mux.HandleFunc(`/`, func(writer http.ResponseWriter, request *http.Request) {
				defer require.NoError(t, request.Body.Close())
				got, err := io.ReadAll(request.Body)
				require.NoError(t, err)
				if !bytes.Equal(got, []byte(tt.body)) {
					t.Errorf("Body changed. Got: %s, Want: %s", got, tt.body)
				}
			})
			h := mw(mux)
			r := httptest.NewRequest(http.MethodPost, `/`, strings.NewReader(tt.body))
			r.Header.Set(headerName, tt.header)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			got := w.Code
			if got != tt.want {
				t.Errorf("HashCheckMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
