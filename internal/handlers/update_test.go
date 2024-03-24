package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kosalnik/metrics/internal/application/server"
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/infra/storage"
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
			storage: storage.NewStorage(),
			method:  http.MethodPost,
			path:    "/update/counter/asdf/3",
			want:    want{statusCode: http.StatusOK},
		},
		{
			name:    "Send gauge",
			storage: storage.NewStorage(),
			method:  http.MethodPost,
			path:    "/update/gauge/asdf/3.0",
			want:    want{statusCode: http.StatusOK},
		},
		{
			name:    "Wrong counter value",
			storage: storage.NewStorage(),
			method:  http.MethodPost,
			path:    "/update/counter/asdf/3.3",
			want:    want{statusCode: http.StatusBadRequest},
		},
		{
			name:    "Wrong gauge value",
			storage: storage.NewStorage(),
			method:  http.MethodPost,
			path:    "/update/gauge/asdf/zxc",
			want:    want{statusCode: http.StatusBadRequest},
		},
		{
			name:    "Wrong type",
			storage: storage.NewStorage(),
			method:  http.MethodPost,
			path:    "/update/zzz/val/1",
			want:    want{statusCode: http.StatusBadRequest},
		},
		{
			name:    "No value",
			storage: storage.NewStorage(),
			method:  http.MethodPost,
			path:    "/update/zzz/val",
			want:    want{statusCode: http.StatusNotFound},
		},
		{
			name:    "No metric name",
			storage: storage.NewStorage(),
			method:  http.MethodPost,
			path:    "/update/zzz",
			want:    want{statusCode: http.StatusNotFound},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := server.NewApp(config.Server{})
			r := app.GetRouter()
			srv := httptest.NewServer(r)
			response, err := srv.Client().Post(srv.URL+tt.path, "text/plain", nil)
			require.NoError(t, err)
			err = response.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}
