// This is agent.
package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"syscall"

	"github.com/kosalnik/metrics/internal/application/client"
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
	cfg := config.NewAgent()
	if err := config.ParseAgentFlags(os.Args, cfg); err != nil {
		panic(err.Error())
	}
	if err := log.InitLogger(cfg.Logger.Level); err != nil {
		panic(err.Error())
	}
	if cfg.Profiling.Enabled {
		go func() {
			if err := http.ListenAndServe(":18080", nil); err != nil {
				panic(err)
			}
		}()
	}
	ctx := context.Background()
	app := client.NewClient(ctx, cfg)
	graceful.
		NewManager(app).
		Notify(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM).
		Run(ctx)
}
