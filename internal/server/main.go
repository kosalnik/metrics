package server

import (
	"github.com/kosalnik/metrics/internal/handler"
	"github.com/kosalnik/metrics/internal/storage"
	"net/http"
)

type App struct {
	Storage storage.Storage
}

func NewApp() *App {
	return &App{
		Storage: storage.NewStorage(),
	}
}

func (app *App) Serve() error {
	mux := http.NewServeMux()
	mux.Handle(`/update/`, handler.NewUpdateHandler(app.Storage))
	return http.ListenAndServe(`:8080`, mux)
}
