package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kosalnik/metrics/internal/application/server"
	"github.com/kosalnik/metrics/internal/config"
)

func TestGetAllHandler(t *testing.T) {
	app := server.NewApp(config.Server{})
	s := app.Storage
	r := app.GetRouter()
	s.IncCounter("c1", 5)
	s.SetGauge("g1", 13.1)
	srv := httptest.NewServer(r)

	res, err := srv.Client().Get(srv.URL)
	require.NoError(t, err)
	content, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	err = res.Body.Close()
	require.NoError(t, err)
	expected := "c1 = 5\ng1 = 13.1"
	assert.Equal(t, expected, string(content))
}
