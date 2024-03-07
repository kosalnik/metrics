package main

import (
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/server"
)

func main() {
	cfg := config.NewConfig()
	parseFlags(&cfg.Server)
	app := server.NewApp(cfg.Server)
	err := app.Serve()
	if err != nil {
		panic(err)
	}
}
