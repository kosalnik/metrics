package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/kosalnik/metrics/internal/config"
)

func parseFlags(c *config.Agent) {
	flag.StringVar(&c.CollectorAddress, "a", "127.0.0.1:8080", "address server endpoint")
	flag.Int64Var(&c.PollInterval, "p", 2, "Pool interval (seconds)")
	flag.Int64Var(&c.ReportInterval, "r", 10, "Report interval (seconds)")
	flag.Int64Var(&c.RateLimit, "l", 1, "Rate limit")
	flag.StringVar(&c.Hash.Key, "k", "", "SHA256 Key")
	flag.Parse()

	var err error
	if v := os.Getenv("PROFILING"); v != "" {
		c.Profiling.Enabled, err = strconv.ParseBool(v)
		if err != nil {
			panic("PROFILING should be bool, got: " + v)
		}
	}
	if v := os.Getenv("ADDRESS"); v != "" {
		c.CollectorAddress = v
	}
	if v := os.Getenv("REPORT_INTERVAL"); v != "" {
		c.ReportInterval, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic("REPORT_INTERVAL should be Int64, got: " + v)
		}
	}
	if v := os.Getenv("RATE_LIMIT"); v != "" {
		c.RateLimit, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic("RATE_LIMIT should be Int64, got: " + v)
		}
	}
	if v := os.Getenv("POLL_INTERVAL"); v != "" {
		c.PollInterval, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic("POLL_INTERVAL should be Int64, got: " + v)
		}
	}
	if v := os.Getenv("KEY"); v != "" {
		c.Hash.Key = v
	}
}
