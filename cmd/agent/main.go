package main

import (
	"context"

	"github.com/kosalnik/metrics/internal/application/client"
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/infra/logger"
)

func main() {
	cfg := config.NewConfig()
	parseFlags(&cfg.Agent)
	if err := logger.InitLogger(cfg.Agent.Logger); err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	app := client.NewClient(ctx, cfg.Agent)
	app.Run(ctx)
}
