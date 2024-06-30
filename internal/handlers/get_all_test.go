package handlers_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/memstorage"
)

func TestGetAllHandler(t *testing.T) {
	var err error
	s := memstorage.NewMemStorage()
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

func BenchmarkGetAllHandler(t *testing.B) {
	var err error
	s := memstorage.NewMemStorage()
	h := handlers.NewGetAllHandler(s)
	for i := 0; i < 1000; i++ {
		_, err = s.IncCounter(context.Background(), fmt.Sprintf("c%d", i), 5)
		assert.NoError(t, err)
		_, err = s.SetGauge(context.Background(), fmt.Sprintf("g%d", i), 13.1)
		assert.NoError(t, err)
	}
	m := http.NewServeMux()
	m.HandleFunc("/", h)
	srv := httptest.NewServer(m)

	for i := 0; i < t.N; i++ {
		r, _ := srv.Client().Get(srv.URL)
		r.Body.Close()
	}
}
