// This is metrics collector server.
package main

import (
	"context"

	"github.com/kosalnik/metrics/internal/application/server"
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/logger"
)

func main() {
	cfg := config.NewConfig()
	parseFlags(&cfg.Server)
	app := server.NewApp(cfg.Server)
	if err := logger.InitLogger(cfg.Server.Logger.Level); err != nil {
		panic(err.Error())
	}
	err := app.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
