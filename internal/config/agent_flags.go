package config

import (
	"flag"
	"os"
	"strconv"
)

const (
	defaultCollectorAddress = "127.0.0.1:8080"
	defaultPollInterval     = 2
	defaultReportInterval   = 10
	defaultRateLimit        = 1
)

func ParseAgentFlags(args []string, c *Agent) {
	fs := flag.NewFlagSet(args[0], flag.PanicOnError)
	fs.SetOutput(os.Stdout)
	fs.StringVar(&c.CollectorAddress, "a", defaultCollectorAddress, "address server endpoint")
	fs.Int64Var(&c.PollInterval, "p", defaultPollInterval, "Pool interval (seconds)")
	fs.Int64Var(&c.ReportInterval, "r", defaultReportInterval, "Report interval (seconds)")
	fs.Int64Var(&c.RateLimit, "l", defaultRateLimit, "Rate limit")
	fs.StringVar(&c.Hash.Key, "k", "", "SHA256 Key")
	if err := fs.Parse(args[1:]); err != nil {
		panic(err.Error())
	}

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
