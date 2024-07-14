package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func ParseAgentFlags(args []string, c *Agent) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.StringVar(&c.CollectorAddress, "a", c.CollectorAddress, "address server endpoint")
	fs.Int64Var(&c.PollInterval, "p", c.PollInterval, "Pool interval (seconds)")
	fs.Int64Var(&c.ReportInterval, "r", c.ReportInterval, "Report interval (seconds)")
	fs.Int64Var(&c.RateLimit, "l", c.RateLimit, "Rate limit")
	fs.StringVar(&c.Hash.Key, "k", c.Hash.Key, "SHA256 Key")
	publicKeyFile := fs.String("crypto-key", "", "Public Key")
	_ = fs.String("config", "", "Config file")
	_ = fs.String("c", "", "Config file (shorthand)")

	var err error
	if err = loadFromConfigFile(fs, c); err != nil {
		return err
	}

	if err = fs.Parse(args[1:]); err != nil {
		return err
	}

	if v := os.Getenv("PROFILING"); v != "" {
		c.Profiling.Enabled, err = strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("PROFILING should be bool, got: %s : %w", v, err)
		}
	}
	if v := os.Getenv("ADDRESS"); v != "" {
		c.CollectorAddress = v
	}
	if v := os.Getenv("REPORT_INTERVAL"); v != "" {
		c.ReportInterval, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("REPORT_INTERVAL should be Int64, got: %s : %w", v, err)
		}
	}
	if v := os.Getenv("RATE_LIMIT"); v != "" {
		c.RateLimit, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("RATE_LIMIT should be Int64, got: %s : %w", v, err)
		}
	}
	if v := os.Getenv("POLL_INTERVAL"); v != "" {
		c.PollInterval, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("POLL_INTERVAL should be Int64, got: %s : %w", v, err)
		}
	}
	if v := os.Getenv("KEY"); v != "" {
		c.Hash.Key = v
	}
	if v := os.Getenv("CRYPTO_KEY"); v != "" {
		publicKeyFile = &v
	}

	if publicKeyFile != nil {
		publicKeyPEM, err := os.ReadFile(*publicKeyFile)
		if err != nil {
			return fmt.Errorf("fail to read key: %w", err)
		}
		publicKeyBlock, _ := pem.Decode(publicKeyPEM)
		publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
		if err != nil {
			return fmt.Errorf("fail to parse key: %w", err)
		}
		c.PublicKey = publicKey.(*rsa.PublicKey)
	}
	return nil
}
