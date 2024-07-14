package config

import (
	"encoding/json"
	"io"
)

func LoadConfigFromFile(f io.Reader, c any) error {
	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, c); err != nil {
		return err
	}
	return nil
}
