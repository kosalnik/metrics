package main

import (
	"flag"
	"github.com/kosalnik/metrics/internal/config"
	"os"
)

func parseFlags(c *config.ServerConfig) {
	flag.StringVar(&c.Address, "a", ":8080", "server endpoint (ip:port)")
	flag.Parse()
	if v := os.Getenv("ADDRESS"); v != "" {
		c.Address = v
	}
}
