package main

import (
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/server"
	"github.com/kosalnik/metrics/internal/service/logger"
)

func main() {
	cfg := config.NewConfig()
	parseFlags(&cfg.Server)
	app := server.NewApp(cfg.Server)
	if err := logger.InitLogger(); err != nil {
		panic(err.Error())
	}
	err := app.Serve()
	if err != nil {
		panic(err)
	}
}
