package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/memstorage"

	"github.com/kosalnik/metrics/internal/storage"
)

func TestUpdateHandler_Handle(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		storage storage.Storage
		method  string
		path    string
		want    want
	}{
		{
			name:    "Send counter",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			path:    "/update/counter/asdf/3",
			want:    want{statusCode: http.StatusOK},
		},
		{
			name:    "Send gauge",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			path:    "/update/gauge/asdf/3.0",
			want:    want{statusCode: http.StatusOK},
		},
		{
			name:    "Wrong counter value",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			path:    "/update/counter/asdf/3.3",
			want:    want{statusCode: http.StatusBadRequest},
		},
		{
			name:    "Wrong gauge value",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			path:    "/update/gauge/asdf/zxc",
			want:    want{statusCode: http.StatusBadRequest},
		},
		{
			name:    "Wrong type",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			path:    "/update/zzz/val/1",
			want:    want{statusCode: http.StatusBadRequest},
		},
		{
			name:    "No value",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			path:    "/update/zzz/val",
			want:    want{statusCode: http.StatusNotFound},
		},
		{
			name:    "No metric name",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			path:    "/update/zzz",
			want:    want{statusCode: http.StatusNotFound},
		},
	}

	s := memstorage.NewMemStorage()
	h := handlers.NewUpdateHandler(s)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", h)
	srv := httptest.NewServer(r)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := srv.Client().Post(srv.URL+tt.path, "text/plain", nil)
			require.NoError(t, err)
			err = response.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}

func TestRestUpdateHandler_Handle(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		storage storage.Storage
		method  string
		req     string
		want    want
	}{
		{
			name:    "Send counter",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			req:     `{"type":"counter","id":"asdf","delta":3}`,
			want:    want{statusCode: http.StatusOK},
		},
		{
			name:    "Send gauge",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			req:     `{"type":"gauge","id":"asdf","value":3.0}`,
			want:    want{statusCode: http.StatusOK},
		},
		{
			name:    "Wrong counter value",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			req:     `{"type":"counter","id":"asdf","delta":3.3}`,
			want:    want{statusCode: http.StatusBadRequest},
		},
		{
			name:    "Wrong gauge value",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			req:     `{"type":"gauge","id":"asdf","value":"zxc"}`,
			want:    want{statusCode: http.StatusBadRequest},
		},
		{
			name:    "Wrong type",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			req:     `{"type":"zzz","id":"val","delta":1}`,
			want:    want{statusCode: http.StatusNotFound},
		},
		{
			name:    "No value",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			req:     `{"type":"zzz","id":"val"}`,
			want:    want{statusCode: http.StatusNotFound},
		},
		{
			name:    "No metric name",
			storage: memstorage.NewMemStorage(),
			method:  http.MethodPost,
			req:     `{"type":"zzz","id":""}`,
			want:    want{statusCode: http.StatusNotFound},
		},
	}

	s := memstorage.NewMemStorage()
	h := handlers.NewRestUpdateHandler(s)

	r := chi.NewRouter()
	r.Post("/", h)
	srv := httptest.NewServer(r)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := srv.Client().Post(srv.URL, "application/json", strings.NewReader(tt.req))
			require.NoError(t, err)
			require.NoError(t, response.Body.Close())
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}
