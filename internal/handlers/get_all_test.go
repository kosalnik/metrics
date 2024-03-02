package handlers_test

import (
	"github.com/kosalnik/metrics/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllHandler(t *testing.T) {
	app := server.NewApp()
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
	expected := "g1 = 13.1\nc1 = 5\n"
	assert.Equal(t, expected, string(content))
}
