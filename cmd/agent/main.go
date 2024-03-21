package main

import (
	"github.com/kosalnik/metrics/internal/client"
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/service/logger"
)

func main() {
	cfg := config.NewConfig()
	parseFlags(&cfg.Agent)
	if err := logger.InitLogger(cfg.Agent.Logger); err != nil {
		panic(err.Error())
	}
	app := client.NewClient(cfg.Agent)
	app.Run()
}
