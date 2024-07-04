package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPingHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	pingerMock := mock.NewMockPinger(ctrl)
	pingerMock.EXPECT().Ping(gomock.Any()).Return(nil)
	h := handlers.NewPingHandler(pingerMock)
	r := chi.NewRouter()
	r.Get("/", h)
	srv := httptest.NewServer(r)
	resp, err := srv.Client().Get(srv.URL)
	require.NoError(t, err)
	defer require.NoError(t, resp.Body.Close())
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestNewPingHandler_FailureByTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	pingerMock := mock.NewMockPinger(ctrl)
	pingerMock.EXPECT().Ping(gomock.Any()).DoAndReturn(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * 2):
			return nil // Прикидываемся что всё хорошо. В тесте контекст должен ошибку вернуть, а не мы
		}
	})
	h := handlers.NewPingHandler(pingerMock)
	r := chi.NewRouter()
	r.Get("/", h)
	srv := httptest.NewServer(r)
	resp, err := srv.Client().Get(srv.URL)
	require.NoError(t, err)
	defer require.NoError(t, resp.Body.Close())
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestNewPingHandler_FailureByDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	pingerMock := mock.NewMockPinger(ctrl)
	pingerMock.EXPECT().Ping(gomock.Any()).Return(errors.New("some error"))
	h := handlers.NewPingHandler(pingerMock)
	r := chi.NewRouter()
	r.Get("/", h)
	srv := httptest.NewServer(r)
	resp, err := srv.Client().Get(srv.URL)
	require.NoError(t, err)
	defer require.NoError(t, resp.Body.Close())
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
