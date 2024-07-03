// This is metrics collector server.
package main

import (
	"context"
	"os"

	"github.com/kosalnik/metrics/internal/application/server"
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/logger"
	"github.com/kosalnik/metrics/internal/version"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	version.Build{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}.Print(os.Stdout)
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
