package server

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/infra/backup"
	"github.com/kosalnik/metrics/internal/infra/logger"
	"github.com/kosalnik/metrics/internal/infra/memstorage"
	"github.com/kosalnik/metrics/internal/infra/postgres"
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

func (app *App) Run(ctx context.Context) error {
	if err := app.initStorage(ctx); err != nil {
		return err
	}
	if err := app.initBackup(ctx); err != nil {
		return err
	}
	defer func() {
		if err := app.Storage.Close(); err != nil {
			logger.Logger.WithError(err).Errorf("unable to close storage")
		}
	}()

	logger.Logger.Info("Listen " + app.config.Address)

	return http.ListenAndServe(app.config.Address, app.GetRouter())
}

func (app *App) initStorage(ctx context.Context) error {
	if app.config.DB.DSN == "" {
		app.Storage = memstorage.NewMemStorage()
	} else {
		return app.initDB(ctx)
	}

	return nil
}

func (app *App) initDB(ctx context.Context) error {
	db, err := postgres.NewDB(ctx, app.config.DB)
	if err != nil {
		return err
	}
	app.Storage = db

	return nil
}

func (app *App) initBackup(ctx context.Context) error {
	var err error
	bm, err := backup.NewBackupManager(app.Storage, app.config.Backup)
	if err != nil {
		return err
	}
	if err = bm.Recover(ctx); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	if err = bm.ScheduleBackup(ctx); err != nil {
		return err
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
		r.With(requireJSONMw).Post("/updates/", handlers.NewUpdateBatchHandler(app.Storage))
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
