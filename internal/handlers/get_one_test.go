package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/memstorage"
)

func TestGetHandler(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		content string
		status  int
	}{
		{"valid gauge", "/value/gauge/g1", "13.1", http.StatusOK},
		{"valid counter", "/value/counter/c1", "5", http.StatusOK},
		{"invalid gauge", "/value/gauge/unknownGauge", "404 page not found\n", http.StatusNotFound},
		{"invalid counter", "/value/counter/unknownCounter", "404 page not found\n", http.StatusNotFound},
		{"invalid metric type", "/value/unk/u3", "Not Found\n", http.StatusNotFound},
	}

	var err error
	s := memstorage.NewMemStorage()
	h := handlers.NewGetHandler(s)
	_, err = s.IncCounter(context.Background(), "c1", 5)
	assert.NoError(t, err)
	_, err = s.SetGauge(context.Background(), "g1", 13.1)
	assert.NoError(t, err)

	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", h)
	srv := httptest.NewServer(r)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := srv.Client().Get(srv.URL + tt.path)
			require.NoError(t, err)
			content, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.status, res.StatusCode)
			err = res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.content, string(content))
		})
	}
}
