package crypt_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/infra/crypt"
)

func TestVerifyHashInterceptor(t *testing.T) {
	tests := map[string]struct {
		name       string
		cfg        config.Hash
		body       string
		wantHeader string
	}{
		"Success. Empty Key": {
			cfg:        config.Hash{Key: ""},
			body:       `asdf`,
			wantHeader: ``,
		},
		"Success. With Key": {
			cfg:        config.Hash{Key: "Secret"},
			body:       `Hello`,
			wantHeader: `6cd4180752f6880f553f2c73c89efe222166924f7bb9707c6240e4b88de77122`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(`/`, func(writer http.ResponseWriter, request *http.Request) {
				require.Equal(t, tt.wantHeader, request.Header.Get(headerName))
			})
			srv := httptest.NewServer(mux)

			client := srv.Client()
			oldTransport := srv.Client().Transport
			client.Transport = crypt.VerifyHashInterceptor(tt.cfg, oldTransport)

			res, err := client.Post(srv.URL, "application/json", strings.NewReader(tt.body))
			require.NoError(t, err)
			require.NoError(t, res.Body.Close())
			require.Equal(t, http.StatusOK, res.StatusCode)
		})
	}
}
