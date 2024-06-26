package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/kosalnik/metrics/internal/config"
)

func parseFlags(c *config.Server) {
	flag.StringVar(&c.Address, "a", ":8080", "server endpoint (ip:port)")
	flag.IntVar(&c.Backup.StoreInterval, "i", 300, "Store interval")
	flag.StringVar(&c.Backup.FileStoragePath, "f", "/tmp/metrics-db.json", "File storage path")
	flag.BoolVar(&c.Backup.Restore, "r", true, "Restore storage before start")
	flag.StringVar(&c.DB.DSN, "d", "", "Database DSN")
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
		c.Address = v
	}
	if v := os.Getenv("STORE_INTERVAL"); v != "" {
		c.Backup.StoreInterval, err = strconv.Atoi(v)
		if err != nil {
			panic("wrong env STORE_INTERVAL")
		}
	}
	if v := os.Getenv("FILE_STORAGE_PATH"); v != "" {
		c.Backup.FileStoragePath = v
	}
	if v := os.Getenv("RESTORE"); v != "" {
		c.Backup.Restore, err = strconv.ParseBool(v)
		if err != nil {
			panic("wrong env RESTORE")
		}
	}
	if v := os.Getenv("DATABASE_DSN"); v != "" {
		c.DB.DSN = v
	}
	if v := os.Getenv("KEY"); v != "" {
		c.Hash.Key = v
	}
}
