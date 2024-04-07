package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kosalnik/metrics/internal/infra/postgres"
	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/infra/storage"
)

type App struct {
	Storage storage.Storage
	config  config.Server
}

func NewApp(cfg config.Server) *App {

	return &App{
		config: cfg,
	}
}

func (app *App) Run() error {
	ctx := context.Background()
	if err := app.initStorage(ctx, app.config); err != nil {
		return err
	}
	defer func() {
		if err := app.Storage.Close(); err != nil {
			logrus.WithError(err).Errorf("unable to close storage")
		}
	}()

	logrus.Info("Listen " + app.config.Address)

	return http.ListenAndServe(app.config.Address, app.GetRouter())
}

func (app *App) initStorage(ctx context.Context, cfg config.Server) error {
	if cfg.DB.DSN == "" {
		storeInterval := time.Second * time.Duration(cfg.StoreInterval)
		app.Storage = storage.NewMemStorage(&storeInterval, &cfg.FileStoragePath)
		if err := app.initBackup(ctx); err != nil {
			return err
		}
	} else {
		return app.initDB(ctx, cfg.DB)
	}

	return nil
}

func (app *App) initDB(ctx context.Context, cfg config.DB) error {
	db, err := postgres.NewDB(ctx, cfg)
	if err != nil {
		return err
	}
	app.Storage = db

	return nil
}

func (app *App) initBackup(ctx context.Context) error {
	if app.config.FileStoragePath == "" {
		return nil
	}
	if app.config.Restore {
		if err := app.Storage.Recover(ctx, app.config.FileStoragePath); err != nil {
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
		r.Get("/ping", handlers.NewPingHandler(app.Storage))
	})
	return r
}
