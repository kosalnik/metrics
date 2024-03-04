package config

type Config struct {
	Agent  Agent
	Server Server
}

type Agent struct {
	// Адрес сервера, куда клиент будет посылать метрики
	CollectorAddress string
	// Время между сборами метрик
	PoolInterval int64
	// Время между отправками метрик на сервер
	ReportInterval int64
}

type Server struct {
	// ip:host, которые слушает сервер
	Address string
}

func NewConfig() *Config {
	return &Config{
		Agent: Agent{
			CollectorAddress: "127.0.0.1:8080",
			PoolInterval:     2,
			ReportInterval:   10,
		},
		Server: Server{
			Address: ":8080",
		},
	}
}
