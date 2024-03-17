package config

type Config struct {
	Agent  Agent
	Server Server
}

type Agent struct {
	Logger Logger
	// Адрес сервера, куда клиент будет посылать метрики
	CollectorAddress string
	// Время между сборами метрик
	PoolInterval int64
	// Время между отправками метрик на сервер
	ReportInterval int64
}

type Server struct {
	Logger Logger
	// ip:host, которые слушает сервер
	Address string
}

type Logger struct {
	Level string
}

func NewConfig() *Config {
	return &Config{
		Agent: Agent{
			Logger:           Logger{Level: "info"},
			CollectorAddress: "127.0.0.1:8080",
			PoolInterval:     2,
			ReportInterval:   10,
		},
		Server: Server{
			Logger:  Logger{Level: "info"},
			Address: ":8080",
		},
	}
}
