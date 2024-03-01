package application

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
		Storage: storage.GetStorage(),
	}
}

func (app *App) Serve() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/gauge/`, handler.HandleUpdateGauge(app.Storage))
	mux.HandleFunc(`/update/counter/`, handler.HandleUpdateCounter(app.Storage))
	return http.ListenAndServe(`:8080`, mux)
}
