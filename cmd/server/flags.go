package main

import (
	"flag"
	"github.com/kosalnik/metrics/internal/config"
)

func parseFlags(c config.ServerConfig) {
	flag.StringVar(&c.Address, "a", ":8080", "server endpoint (ip:port)")
}
