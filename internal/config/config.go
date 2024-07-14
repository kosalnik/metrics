// Package config contains config.
// Ну что тут сказать. В этом пакете находятся структуры конфигураций компонентов системы.
package config

import (
	"crypto/rsa"

	"github.com/kosalnik/metrics/internal/backup"
	"github.com/kosalnik/metrics/internal/crypt"
	"github.com/kosalnik/metrics/internal/log"
)

type Config struct {
	Server Server
	Agent  Agent
}

type Agent struct {
	CollectorAddress string
	Hash             crypt.Config
	Logger           log.Config
	PollInterval     int64
	ReportInterval   int64
	RateLimit        int64
	Profiling        Profiling
	PublicKey        *rsa.PublicKey
}

type Server struct {
	Logger     log.Config
	DB         DB
	Hash       crypt.Config
	Address    string
	Backup     backup.Config
	Profiling  Profiling
	PrivateKey *rsa.PrivateKey
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
			Logger:           log.Config{Level: "info"},
			CollectorAddress: "127.0.0.1:8080",
			PollInterval:     2,
			ReportInterval:   10,
			Hash:             crypt.Config{Key: ""},
			RateLimit:        1,
		},
		Server: Server{
			Profiling: Profiling{},
			Logger:    log.Config{Level: "info"},
			Address:   ":8080",
			Backup:    backup.Config{},
			Hash:      crypt.Config{Key: ""},
		},
	}
}

func NewAgent() *Agent {
	return &Agent{
		Profiling:        Profiling{},
		Logger:           log.Config{Level: "info"},
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
		Logger:    log.Config{Level: "info"},
		Address:   ":8080",
		Backup:    backup.Config{},
		Hash:      crypt.Config{Key: ""},
	}
}
