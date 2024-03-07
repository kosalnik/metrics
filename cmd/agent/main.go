package main

import (
	"github.com/kosalnik/metrics/internal/client"
	"github.com/kosalnik/metrics/internal/config"
)

func main() {
	cfg := config.NewConfig()
	parseFlags(&cfg.Agent)
	app := client.NewClient(cfg.Agent)
	app.Run()
}
