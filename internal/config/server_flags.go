package config

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	defaultAddress               = ":8080"
	defaultStoreInterval         = 300
	defaultBackupFileStoragePath = "/tmp/metrics-db.json"
)

func ParseServerFlags(args []string, c *Server) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.StringVar(&c.Address, "a", defaultAddress, "server endpoint (ip:port)")
	fs.IntVar(&c.Backup.StoreInterval, "i", defaultStoreInterval, "Store interval")
	fs.StringVar(&c.Backup.FileStoragePath, "f", defaultBackupFileStoragePath, "File storage path")
	fs.BoolVar(&c.Backup.Restore, "r", true, "Restore storage before start")
	fs.StringVar(&c.DB.DSN, "d", "", "Database DSN")
	fs.StringVar(&c.Hash.Key, "k", "", "SHA256 Key")
	privateKeyFile := fs.String("crypto-key", "", "Public Key")
	_ = fs.String("config", "", "Config file")
	_ = fs.String("c", "", "Config file (shorthand)")

	var err error
	if err = loadFromConfigFile(fs, c); err != nil {
		return err
	}

	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("fail to parse flags: %w", err)
	}
	if v := os.Getenv("PROFILING"); v != "" {
		c.Profiling.Enabled, err = strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("PROFILING should be bool, got: %s : %w", v, err)
		}
	}
	if v := os.Getenv("ADDRESS"); v != "" {
		c.Address = v
	}
	if v := os.Getenv("STORE_INTERVAL"); v != "" {
		c.Backup.StoreInterval, err = strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("wrong env STORE_INTERVAL: %w", err)
		}
	}
	if v := os.Getenv("FILE_STORAGE_PATH"); v != "" {
		c.Backup.FileStoragePath = v
	}
	if v := os.Getenv("RESTORE"); v != "" {
		c.Backup.Restore, err = strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("wrong env RESTORE: %w", err)
		}
	}
	if v := os.Getenv("DATABASE_DSN"); v != "" {
		c.DB.DSN = v
	}
	if v := os.Getenv("KEY"); v != "" {
		c.Hash.Key = v
	}
	if v := os.Getenv("CRYPTO_KEY"); v != "" {
		privateKeyFile = &v
	}

	if privateKeyFile != nil && *privateKeyFile != "" {
		privateKeyPEM, err := os.ReadFile(*privateKeyFile)
		if err != nil {
			return fmt.Errorf("fail to read key: %w", err)
		}
		keyBlock, _ := pem.Decode(privateKeyPEM)
		privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if err != nil {
			return fmt.Errorf("fail to parse key: %w", err)
		}
		c.PrivateKey = privateKey
	}
	return nil
}
