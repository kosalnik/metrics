package main

import (
	"flag"
	"os"

	"github.com/kosalnik/metrics/internal/config"
)

func parseFlags(c *config.Server) {
	flag.StringVar(&c.Address, "a", ":8080", "server endpoint (ip:port)")
	flag.Parse()
	if v := os.Getenv("ADDRESS"); v != "" {
		c.Address = v
	}
}
