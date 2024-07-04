package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUpdateBatchHandler(t *testing.T) {
	cases := map[string]struct {
		req            string
		wantStatusCode int
	}{
		"Success": {
			`[{"id":"a","type":"counter","delta":3},{"id":"b","type":"gauge","value":3.14}]`,
			http.StatusOK,
		},
		"Wrong json": {
			`{"id":"a","type":"counter","delta":3}`,
			http.StatusBadRequest,
		},
	}
	s := memstorage.NewMemStorage()
	h := handlers.NewUpdateBatchHandler(s)
	r := chi.NewRouter()
	r.Post("/", h)
	srv := httptest.NewServer(r)
	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			resp, err := srv.Client().Post(srv.URL, "application/json", strings.NewReader(tt.req))
			require.NoError(t, err)
			defer require.NoError(t, resp.Body.Close())
			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
		})
	}
}
