package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/infra/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllHandler(t *testing.T) {
	var err error
	s := storage.NewMemStorage(nil, nil)
	h := handlers.NewGetAllHandler(s)
	_, err = s.IncCounter(context.Background(), "c1", 5)
	assert.NoError(t, err)
	_, err = s.SetGauge(context.Background(), "g1", 13.1)
	assert.NoError(t, err)

	m := http.NewServeMux()
	m.HandleFunc("/", h)
	srv := httptest.NewServer(m)

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
