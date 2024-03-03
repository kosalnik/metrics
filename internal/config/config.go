package config

import "time"

type Config struct {
	Agent  AgentConfig
	Server ServerConfig
}

type AgentConfig struct {
	// Адрес сервера, куда клиент будет посылать метрики
	CollectorAddress string
	// Время между сборами метрик
	PoolInterval time.Duration
	// Время между отправками метрик на сервер
	ReportInterval time.Duration
}

type ServerConfig struct {
	// ip:host, которые слушает сервер
	Address string
}

func NewConfig() *Config {
	return &Config{
		Agent: AgentConfig{
			CollectorAddress: "http://127.0.0.1:8080",
			PoolInterval:     time.Duration(time.Second * 2),
			ReportInterval:   time.Duration(time.Second * 10),
		},
		Server: ServerConfig{
			Address: ":8080",
		},
	}
}
