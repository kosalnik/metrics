package config

import (
	"flag"
	"os"
	"strconv"
)

const (
	defaultAddress               = ":8080"
	defaultStoreInterval         = 300
	defaultBackupFileStoragePath = "/tmp/metrics-db.json"
)

func ParseServerFlags(args []string, c *Server) {
	fs := flag.NewFlagSet(args[0], flag.PanicOnError)
	fs.SetOutput(os.Stdout)
	fs.StringVar(&c.Address, "a", defaultAddress, "server endpoint (ip:port)")
	fs.IntVar(&c.Backup.StoreInterval, "i", defaultStoreInterval, "Store interval")
	fs.StringVar(&c.Backup.FileStoragePath, "f", defaultBackupFileStoragePath, "File storage path")
	fs.BoolVar(&c.Backup.Restore, "r", true, "Restore storage before start")
	fs.StringVar(&c.DB.DSN, "d", "", "Database DSN")
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
