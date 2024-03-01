package config

import "time"

type Config struct {
	Client ClientConfig
}

type ClientConfig struct {
	CollectorAddress string
	PoolInterval     time.Duration
	ReportInterval   time.Duration
}

func NewConfig() *Config {
	return &Config{
		Client: ClientConfig{
			CollectorAddress: "http://127.0.0.1:8080",
			PoolInterval:     time.Duration(2_000_000_000),
			ReportInterval:   time.Duration(10_000_000_000),
		},
	}
}
