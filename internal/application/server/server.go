package server

import (
	"net/http"
	"strings"

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
		gzipMiddleware,
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
	})
	return r
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
