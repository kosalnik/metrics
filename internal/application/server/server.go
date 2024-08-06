// Package server - пакет с реализацией сервера сбора метрик.
// При старте методом App.Run() пытается восстановить storage из бекапа, затем периодически сбрасывает storage в бекап.
// Запускает инстанс http.Server, который принимает метрики от внешнего сервиса.
package server

import (
	"context"
	"crypto/rand"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	pb "github.com/kosalnik/metrics/pkg/metrics"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/kosalnik/metrics/internal/backup"
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/crypt"
	"github.com/kosalnik/metrics/internal/handlers"
	"github.com/kosalnik/metrics/internal/log"
	"github.com/kosalnik/metrics/internal/memstorage"
	"github.com/kosalnik/metrics/internal/postgres"
	"github.com/kosalnik/metrics/internal/storage"
)

type App struct {
	Storage       storage.Storager
	config        *config.Server
	server        *http.Server
	grpcServer    *grpc.Server
	backupManager *backup.BackupManager
}

func NewApp(cfg *config.Server) *App {
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
			log.Error().Err(err).Msg("unable to close storage")
		}
	}()

	wg := sync.WaitGroup{}

	if app.config.GRPCAddress != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info().Str("addr", app.config.GRPCAddress).Msg("Listen grpc")
			listen, err := net.Listen("tcp", app.config.GRPCAddress)
			if err != nil {
				log.Fatal().Err(err).Msg("new grpc server fails")
			}
			app.grpcServer = grpc.NewServer()
			pb.RegisterMetricsServer(app.grpcServer, &GRPCServer{storage: app.Storage})
			if err := app.grpcServer.Serve(listen); err != nil {
				log.Error().Err(err).Msg("Listen grpc fails")
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.server = &http.Server{
			Addr:    app.config.Address,
			Handler: app.GetRouter(),
		}

		log.Info().Str("address", app.config.Address).Msg("Listen")
		if err := app.server.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("Listen http fails")
		}
	}()
	wg.Wait()
	return nil
}

func (app *App) Shutdown(ctx context.Context) {
	log.Info().Msg(`Shutdown start`)
	g := errgroup.Group{}
	g.Go(func() (err error) {
		log.Info().Msg(`Shutdown "server.App" start`)
		defer func() {
			if err != nil {
				log.Error().Err(err).Msg(`Shutdown "server.App" error`)
			} else {
				log.Info().Msg(`Shutdown "server.App" completed`)
			}
		}()
		return app.server.Shutdown(ctx)
	})
	g.Go(func() (err error) {
		log.Info().Msg(`Shutdown "backupManager" start`)
		defer func() {
			if err != nil {
				log.Error().Err(err).Msg(`Shutdown "backupManager" error`)
			} else {
				log.Info().Msg(`Shutdown "backupManager" completed`)
			}
		}()
		err = app.backupManager.Store(ctx)
		return
	})
	if app.grpcServer != nil {
		g.Go(func() (_ error) {
			log.Info().Msg(`Shutdown "grpc" start`)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("Shutdown error")
	}
	log.Info().Msg(`Shutdown completed`)
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
	db, err := postgres.NewConn(app.config.DB.DSN)
	if err != nil {
		return err
	}
	dbs, err := postgres.NewDBStorage(db)
	if err != nil {
		return err
	}
	if err := dbs.CreateTablesIfNotExist(ctx); err != nil {
		return err
	}
	app.Storage = dbs
	return nil
}

// ScheduleBackup - запустить автоматический бекап по расписанию.
// Будет скидывать бекап на диск через равные промежутки времени.
func (app *App) initBackup(ctx context.Context) error {
	var err error
	app.backupManager, err = backup.NewBackupManager(app.Storage, app.config.Backup)
	if err != nil {
		return err
	}
	if err = app.backupManager.Recover(ctx); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	log.Info().Msg("schedule backup")
	go app.backupManager.BackupLoop(ctx)

	return nil
}

func (app *App) GetRouter() chi.Router {
	r := chi.NewRouter()

	if app.config.PrivateKey != nil {
		r.Use(crypt.CipherMiddleware(crypt.NewDecoder(app.config.PrivateKey, rand.Reader)))
	}

	if app.config.TrustedSubnet != "" {
		log.Info().Str("CIDR", app.config.TrustedSubnet).Msg("Protect with trusted subnet")
		_, subnet, err := net.ParseCIDR(app.config.TrustedSubnet)
		if err != nil {
			panic(err)
		}
		r.Use(TrustedClientMiddleware(subnet))
	}
	r.Use(
		middleware.Compress(1, "application/json", "text/html"),
		middleware.Logger,
		middleware.Recoverer,
		crypt.HashCheckMiddleware(app.config.Hash),
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
		if app.config.Profiling.Enabled {
			r.Mount("/profiler", middleware.Profiler())
		}
	})
	return r
}
