package main

import (
	"github.com/kosalnik/metrics/internal/client"
	"github.com/kosalnik/metrics/internal/config"
)

func main() {
	cfg := config.NewConfig()
	app := client.NewClient(cfg.Client)
	app.Run()
}
