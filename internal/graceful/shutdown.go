package graceful

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/kosalnik/metrics/internal/log"
	"golang.org/x/sync/errgroup"
)

type Shutdowner interface {
	Shutdown(ctx context.Context)
}

type Runner interface {
	Run(ctx context.Context) error
}

type ShutdownRunner interface {
	Shutdowner
	Runner
}

type Manager struct {
	shutdownTimeout time.Duration
	services        []ShutdownRunner
	signals         []os.Signal
}

const defaultShutdownTimeout = time.Second * 5

func NewManager(services ...ShutdownRunner) *Manager {
	return &Manager{
		services:        services,
		shutdownTimeout: defaultShutdownTimeout,
	}
}

func (m *Manager) Notify(signals ...os.Signal) *Manager {
	m.signals = signals
	return m
}

func (m *Manager) Run(ctx context.Context) {
	log.Info().Any("signals", m.signals).Msg("Graceful shutdown subscribe on signals")
	ctxNotify, stop := signal.NotifyContext(ctx, m.signals...)
	defer stop()

	g := new(errgroup.Group)
	for i := range m.services {
		g.Go(func() error {
			return m.runService(ctxNotify, m.services[i])
		})
	}

	go func() {
		if err := g.Wait(); err != nil {
			stop()
		}
	}()

	<-ctxNotify.Done()

	log.Info().Err(ctxNotify.Err()).Msg("Start Graceful shutdown")
	wg := sync.WaitGroup{}
	for i := range m.services {
		wg.Add(1)
		go func(service Shutdowner) {
			ctxTimeout, cancel := context.WithTimeout(context.Background(), m.shutdownTimeout)
			defer cancel()
			log.Info().Type("service", service).Msg("Service shutdown start")
			service.Shutdown(ctxTimeout)
			if errors.Is(ctxTimeout.Err(), context.DeadlineExceeded) {
				log.Error().Msg("Service shutdown fails. Deadline exceed")
			} else {
				log.Info().Type("service", service).Msg("Service shutdown finish")
			}
			wg.Done()
		}(m.services[i])
	}
	log.Info().Msg("Graceful shutdown in progress")
	wg.Wait()
	log.Info().Msg("Graceful shutdown completed")
}

func (m *Manager) runService(ctx context.Context, s Runner) error {
	if err := s.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("Service error")
		return err
	}
	return nil
}
