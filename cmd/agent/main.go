package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"

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
	if cfg.Agent.Profiling.Enabled {
		go func() {
			if err := http.ListenAndServe(":18080", nil); err != nil {
				panic(err)
			}
		}()
	}
	ctx := context.Background()
	app := client.NewClient(ctx, cfg.Agent)
	app.Run(ctx)

}
