package server

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/infra/storage"
)

type App struct {
	Storage storage.Storage
	config  config.Server
	db      *sql.DB
}

func NewApp(cfg config.Server) *App {
	storeInterval := time.Second * time.Duration(cfg.StoreInterval)
	return &App{
		Storage: storage.NewStorage(&storeInterval, &cfg.FileStoragePath),
		config:  cfg,
	}
}

func (app *App) Run() error {
	if err := app.initBackup(); err != nil {
		return err
	}

	if err := app.initDB(app.config.Db); err != nil {
		return err
	}
	defer func() {
		if err := app.db.Close(); err != nil {
			logrus.WithError(err).Errorf("unable to close db")
		}
	}()

	logrus.Info("Listen " + app.config.Address)

	return http.ListenAndServe(app.config.Address, app.GetRouter())
}

func (app *App) initDB(cfg config.Db) error {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return err
	}
	app.db = db

	return nil
}

func (app *App) initBackup() error {
	if app.config.FileStoragePath == "" {
		return nil
	}
	if app.config.Restore {
		if err := app.Storage.Recover(app.config.FileStoragePath); err != nil {
			if errors.Is(os.ErrNotExist, err) {
				return err
			}
		}
	}
	return nil
}

func (app *App) GetRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(
		middleware.Compress(1, "application/json", "text/html"),
		//gzipMiddleware,
		middleware.Logger,
		middleware.Recoverer,
	)
	requireJSONMw := middleware.AllowContentType("application/json")
	r.Route("/", func(r chi.Router) {
		r.With(requireJSONMw).Get("/", handlers.NewGetAllHandler(app.Storage))
		r.Route("/update", func(r chi.Router) {
			r.With(requireJSONMw).Post("/", handlers.NewRestUpdateHandler(app.Storage))
			r.Post("/{type}/{name}/{value}", handlers.NewUpdateHandler(app.Storage))
		})
		r.Route("/value", func(r chi.Router) {
			r.With(requireJSONMw).Post("/", handlers.NewRestGetHandler(app.Storage))
			r.Get("/{type}/{name}", handlers.NewGetHandler(app.Storage))
		})
		r.Get("/ping", handlers.NewPingHandler(app.db))
	})
	return r
}
