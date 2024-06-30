// Package config contains config.
// Ну что тут сказать. В этом пакете находятся структуры конфигураций компонентов системы.
package config

import (
	"github.com/kosalnik/metrics/internal/backup"
	"github.com/kosalnik/metrics/internal/crypt"
	"github.com/kosalnik/metrics/internal/logger"
)

type Config struct {
	Agent  Agent
	Server Server
}

type Agent struct {
	Profiling        Profiling
	Logger           logger.Config
	CollectorAddress string // Адрес сервера, куда клиент будет посылать метрики
	PollInterval     int64  // Время между сборами метрик
	ReportInterval   int64  // Время между отправками метрик на сервер
	RateLimit        int64  //
	Hash             crypt.Config
}

type Server struct {
	Profiling Profiling
	Logger    logger.Config
	Address   string // ip:host, которые слушает сервер
	Backup    backup.Config
	DB        DB
	Hash      crypt.Config
}

type Profiling struct {
	Enabled bool
}

type DB struct {
	DSN string
}

func NewConfig() *Config {
	return &Config{
		Agent: Agent{
			Profiling:        Profiling{},
			Logger:           logger.Config{Level: "info"},
			CollectorAddress: "127.0.0.1:8080",
			PollInterval:     2,
			ReportInterval:   10,
			Hash:             crypt.Config{Key: ""},
			RateLimit:        1,
		},
		Server: Server{
			Profiling: Profiling{},
			Logger:    logger.Config{Level: "info"},
			Address:   ":8080",
			Backup:    backup.Config{},
			Hash:      crypt.Config{Key: ""},
		},
	}
}

func NewAgent() *Agent {
	return &Agent{
		Profiling:        Profiling{},
		Logger:           logger.Config{Level: "info"},
		CollectorAddress: "127.0.0.1:8080",
		PollInterval:     2,
		ReportInterval:   10,
		Hash:             crypt.Config{Key: ""},
		RateLimit:        1,
	}
}

func NewServer() *Server {
	return &Server{
		Profiling: Profiling{},
		Logger:    logger.Config{Level: "info"},
		Address:   ":8080",
		Backup:    backup.Config{},
		Hash:      crypt.Config{Key: ""},
	}
}
