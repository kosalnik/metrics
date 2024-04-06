package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/kosalnik/metrics/internal/config"
)

func parseFlags(c *config.Server) {
	flag.StringVar(&c.Address, "a", ":8080", "server endpoint (ip:port)")
	flag.IntVar(&c.StoreInterval, "i", 300, "Store interval")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "File storage path")
	flag.BoolVar(&c.Restore, "r", true, "Restore storage before start")
	flag.StringVar(&c.DB.DSN, "d", "", "Database DSN")
	flag.Parse()
	var err error
	if v := os.Getenv("ADDRESS"); v != "" {
		c.Address = v
	}
	if v := os.Getenv("STORE_INTERVAL"); v != "" {
		c.StoreInterval, err = strconv.Atoi(v)
		if err != nil {
			panic("wrong env STORE_INTERVAL")
		}
	}
	if v := os.Getenv("FILE_STORAGE_PATH"); v != "" {
		c.FileStoragePath = v
	}
	if v := os.Getenv("RESTORE"); v != "" {
		c.Restore, err = strconv.ParseBool(v)
		if err != nil {
			panic("wrong env RESTORE")
		}
	}
	if v := os.Getenv("DATABASE_DSN"); v != "" {
		c.DB.DSN = v
	}
}
