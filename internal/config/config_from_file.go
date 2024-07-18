package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"

	"github.com/kosalnik/metrics/internal/log"
)

func loadFromJson(f io.Reader, c any) error {
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, c); err != nil {
		return err
	}
	return nil
}

func loadFromConfigFile(fs *flag.FlagSet, c any) error {
	var configFilePath string
	fs.Visit(func(f *flag.Flag) {
		if (f.Name == "config" || f.Name == "c") && f.Value.String() != "" {
			configFilePath = f.Value.String()
		}
	})
	if v := os.Getenv("CONFIG"); v != "" {
		configFilePath = v
	}
	if configFilePath != "" {
		log.Debug().Str("path", configFilePath).Msg("Load config from file")
		fl, err := os.Open(configFilePath)
		if err != nil {
			return err
		}
		if err := loadFromJson(fl, c); err != nil {
			return err
		}
	}
	return nil
}
