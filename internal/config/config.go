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
			PoolInterval:     time.Duration(time.Second * 2),
			ReportInterval:   time.Duration(time.Second * 10),
		},
	}
}
