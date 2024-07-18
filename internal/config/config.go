// Package config contains config.
// Ну что тут сказать. В этом пакете находятся структуры конфигураций компонентов системы.
package config

import (
	"crypto/rsa"
	_ "embed"

	"github.com/kosalnik/metrics/internal/backup"
	"github.com/kosalnik/metrics/internal/crypt"
	"github.com/kosalnik/metrics/internal/log"
)

const (
	defaultServerBind       = ":8080"
	defaultCollectorAddress = "127.0.0.1:8080"
	defaultPollInterval     = 2
	defaultReportInterval   = 10
	defaultRateLimit        = 1
)

type Config struct {
	Server Server
	Agent  Agent
}

type Agent struct {
	CollectorAddress string         `json:"address"`
	PollInterval     int64          `json:"poll_interval"`
	ReportInterval   int64          `json:"report_interval"`
	PublicKey        *rsa.PublicKey `json:"crypto_key"`
	RateLimit        int64          `json:"rate_limit"`
	Hash             crypt.Config
	Logger           log.Config
	Profiling        Profiling
}

type Server struct {
	Address    string          `json:"address"`
	Backup     backup.Config   `json:"backup"`
	PrivateKey *rsa.PrivateKey `json:"crypto_key"`
	Logger     log.Config
	DB         DB
	Hash       crypt.Config
	Profiling  Profiling
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
			Address:   defaultServerBind,
			Backup:    backup.Config{},
			Hash:      crypt.Config{Key: ""},
		},
	}
}

func NewAgent() *Agent {
	return &Agent{
		Profiling:        Profiling{},
		Logger:           log.Config{Level: "info"},
		CollectorAddress: defaultCollectorAddress,
		PollInterval:     defaultPollInterval,
		ReportInterval:   defaultReportInterval,
		Hash:             crypt.Config{Key: ""},
		RateLimit:        defaultRateLimit,
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
