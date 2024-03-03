package main

import (
	"flag"
	"github.com/kosalnik/metrics/internal/config"
)

func parseFlags(c *config.AgentConfig) {
	flag.StringVar(&c.CollectorAddress, "a", "127.0.0.1:8080", "address server endpoint")
	flag.Int64Var(&c.PoolInterval, "p", 2, "Pool interval (seconds)")
	flag.Int64Var(&c.ReportInterval, "r", 10, "Report interval (seconds)")
	flag.Parse()
}
