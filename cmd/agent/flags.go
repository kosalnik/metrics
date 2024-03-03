package main

import (
	"flag"
	"github.com/kosalnik/metrics/internal/config"
	"time"
)

func parseFlags(c config.AgentConfig) {
	flag.StringVar(&c.CollectorAddress, "a", "http://127.0.0.1:8080", "Full address server endpoint")
	flag.DurationVar(&c.PoolInterval, "p", time.Second*2, "Pool interval (seconds)")
	flag.DurationVar(&c.ReportInterval, "r", time.Second*10, "Report interval (seconds)")
}
