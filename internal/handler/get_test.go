package handler_test

import (
	"github.com/kosalnik/metrics/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		content string
		status  int
	}{
		{"", "/value/gauge/g1", "13.1", http.StatusOK},
		{"", "/value/counter/c1", "5", http.StatusOK},
		{"", "/value/unk/u3", "Not Found\n", http.StatusNotFound},
	}

	app := server.NewApp()
	s := app.Storage
	r := app.GetRouter()
	s.IncCounter("c1", 5)
	s.SetGauge("g1", 13.1)
	srv := httptest.NewServer(r)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := srv.Client().Get(srv.URL + tt.path)
			require.NoError(t, err)
			content, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.status, res.StatusCode)
			res.Body.Close()
			assert.Equal(t, tt.content, string(content))
		})
	}
}
