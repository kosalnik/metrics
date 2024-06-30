// Package config contains config.
// Ну что тут сказать. В этом пакете находятся структуры конфигураций компонентов системы.
package config

type Config struct {
	Agent  Agent
	Server Server
}

type Agent struct {
	Profiling        Profiling
	Logger           Logger
	CollectorAddress string // Адрес сервера, куда клиент будет посылать метрики
	PollInterval     int64  // Время между сборами метрик
	ReportInterval   int64  // Время между отправками метрик на сервер
	RateLimit        int64  //
	Hash             Hash
}

type Server struct {
	Profiling Profiling
	Logger    Logger
	Address   string // ip:host, которые слушает сервер
	Backup    Backup
	DB        DB
	Hash      Hash
}

type Profiling struct {
	Enabled bool
}

type Hash struct {
	Key string // HASH SHA256 Key
}

type Backup struct {
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

type DB struct {
	DSN string
}

type Logger struct {
	Level string
}

func NewConfig() *Config {
	return &Config{
		Agent: Agent{
			Profiling:        Profiling{},
			Logger:           Logger{Level: "info"},
			CollectorAddress: "127.0.0.1:8080",
			PollInterval:     2,
			ReportInterval:   10,
			Hash:             Hash{Key: ""},
			RateLimit:        1,
		},
		Server: Server{
			Profiling: Profiling{},
			Logger:    Logger{Level: "info"},
			Address:   ":8080",
			Backup:    Backup{},
			Hash:      Hash{Key: ""},
		},
	}
}
