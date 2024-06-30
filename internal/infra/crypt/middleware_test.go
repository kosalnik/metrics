package crypt_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kosalnik/metrics/internal/infra/crypt"
)

func TestHashCheckMiddleware(t *testing.T) {
	tests := map[string]struct {
		cfg    crypt.Config
		body   string
		header string
		want   int
	}{
		"Success. Empty Key": {
			cfg:    crypt.Config{Key: ""},
			body:   `asdf`,
			header: `f553f2c73c8`,
			want:   http.StatusOK,
		},
		"Success. With Key": {
			cfg:    crypt.Config{Key: "Secret"},
			body:   `Hello`,
			header: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			want:   http.StatusOK,
		},
		"Invalid hash. Other Key": {
			cfg:    crypt.Config{Key: "asd"},
			body:   `Hello`,
			header: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			want:   http.StatusBadRequest,
		},
		"Invalid hash. Other body": {
			cfg:    crypt.Config{Key: "Secret"},
			body:   `asd`,
			header: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			want:   http.StatusBadRequest,
		},
		"Invalid hash. Other hash": {
			cfg:    crypt.Config{Key: "Secret"},
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

func BenchmarkHashCheckMiddleware(b *testing.B) {
	tests := map[string]struct {
		cfg    crypt.Config
		body   string
		header string
		want   int
	}{
		"Success. Empty Key": {
			cfg:    crypt.Config{Key: ""},
			body:   `asdf`,
			header: `f553f2c73c8`,
			want:   http.StatusOK,
		},
		"Success. With Key": {
			cfg:    crypt.Config{Key: "Secret"},
			body:   `Hello`,
			header: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			want:   http.StatusOK,
		},
		"Invalid hash. Other Key": {
			cfg:    crypt.Config{Key: "asd"},
			body:   `Hello`,
			header: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			want:   http.StatusBadRequest,
		},
		"Invalid hash. Other body": {
			cfg:    crypt.Config{Key: "Secret"},
			body:   `asd`,
			header: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
			want:   http.StatusBadRequest,
		},
		"Invalid hash. Other hash": {
			cfg:    crypt.Config{Key: "Secret"},
			body:   `Hello`,
			header: `asdf`,
			want:   http.StatusBadRequest,
		},
	}

	b.ResetTimer()

	for name, tt := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				mw := crypt.HashCheckMiddleware(tt.cfg)
				mux := http.NewServeMux()
				mux.HandleFunc(`/`, func(writer http.ResponseWriter, request *http.Request) {
					defer require.NoError(b, request.Body.Close())
					got, err := io.ReadAll(request.Body)
					require.NoError(b, err)
					if !bytes.Equal(got, []byte(tt.body)) {
						b.Errorf("Body changed. Got: %s, Want: %s", got, tt.body)
					}
				})
				h := mw(mux)
				r := httptest.NewRequest(http.MethodPost, `/`, strings.NewReader(tt.body))
				r.Header.Set(headerName, tt.header)
				w := httptest.NewRecorder()

				b.StartTimer()

				h.ServeHTTP(w, r)

				got := w.Code
				if got != tt.want {
					b.Errorf("HashCheckMiddleware() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
