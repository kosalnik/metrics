// This is metrics collector server.
package main

import (
	"context"
	"os"
	"syscall"

	"github.com/kosalnik/metrics/internal/application/server"
	"github.com/kosalnik/metrics/internal/config"
	"github.com/kosalnik/metrics/internal/graceful"
	"github.com/kosalnik/metrics/internal/log"
	"github.com/kosalnik/metrics/internal/version"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	version.VersionInfo{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}.Print(os.Stdout)
	ctx := context.Background()
	cfg := config.NewConfig()
	if err := config.ParseServerFlags(os.Args, &cfg.Server); err != nil {
		panic(err.Error())
	}
	app := server.NewApp(cfg.Server)
	if err := log.InitLogger(cfg.Server.Logger.Level); err != nil {
		panic(err.Error())
	}
	graceful.
		NewManager(app).
		Notify(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM).
		Run(ctx)
}
