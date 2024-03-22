package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/storage"
)

type App struct {
	Storage storage.Storage
	config  config.Server
}

func NewApp(cfg config.Server) *App {
	return &App{
		Storage: storage.NewStorage(),
		config:  cfg,
	}
}

func (app *App) Serve() error {
	logrus.Info("Listen " + app.config.Address)
	return http.ListenAndServe(app.config.Address, app.GetRouter())
}

func (app *App) GetRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)
	requireJSONMw := middleware.AllowContentType("application/json")
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.NewGetAllHandler(app.Storage))
		r.Route("/update", func(r chi.Router) {
			r.With(requireJSONMw).Post("/", handlers.NewRestUpdateHandler(app.Storage))
			r.Post("/{type}/{name}/{value}", handlers.NewUpdateHandler(app.Storage))
		})
		r.Route("/value", func(r chi.Router) {
			r.With(requireJSONMw).Post("/", handlers.NewRestGetHandler(app.Storage))
			r.Get("/{type}/{name}", handlers.NewGetHandler(app.Storage))
		})
	})
	return r
}
