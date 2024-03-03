package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/storage"
	"log"
	"net/http"
)

type App struct {
	Storage storage.Storage
	config  config.ServerConfig
}

func NewApp(cfg config.ServerConfig) *App {
	return &App{
		Storage: storage.NewStorage(),
		config:  cfg,
	}
}

func (app *App) Serve() error {
	log.Println("Listen " + app.config.Address)
	return http.ListenAndServe(app.config.Address, app.GetRouter())
}

func (app *App) GetRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.NewGetAllHandler(app.Storage))
		r.Post("/update/{type}/{name}/{value}", func(writer http.ResponseWriter, request *http.Request) {
			h := handlers.NewUpdateHandler(app.Storage)
			h.Handle(writer, request)
		})
		r.Get("/value/{type}/{name}", handlers.NewGetHandler(app.Storage))
	})
	return r
}
