package handler

import (
	"github.com/kosalnik/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler_ServeHTTP(t *testing.T) {
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
		{
			name:    "Type not specified",
			storage: storage.NewStorage(),
			method:  http.MethodPost,
			path:    "/update",
			want:    want{statusCode: http.StatusNotFound},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UpdateHandler{
				Storage: tt.storage,
			}
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			h.ServeHTTP(w, req)
			res := w.Result()
			err := res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, w.Code)
		})
	}
}
