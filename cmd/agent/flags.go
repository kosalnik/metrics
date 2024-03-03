package main

import (
	"flag"
	"github.com/kosalnik/metrics/internal/config"
	"os"
	"strconv"
)

func parseFlags(c *config.AgentConfig) {
	flag.StringVar(&c.CollectorAddress, "a", "127.0.0.1:8080", "address server endpoint")
	flag.Int64Var(&c.PoolInterval, "p", 2, "Pool interval (seconds)")
	flag.Int64Var(&c.ReportInterval, "r", 10, "Report interval (seconds)")
	flag.Parse()

	var err error
	if v := os.Getenv("ADDRESS"); v != "" {
		c.CollectorAddress = v
	}
	if v := os.Getenv("REPORT_INTERVAL"); v != "" {
		c.ReportInterval, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic("REPORT_INTERVAL should be Int64, got: " + v)
		}
	}
	if v := os.Getenv("POLL_INTERVAL"); v != "" {
		c.PoolInterval, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic("POLL_INTERVAL should be Int64, got: " + v)
		}
	}
}
