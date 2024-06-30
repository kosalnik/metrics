// Package logger initializes logger.
package logger

import (
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

type Config struct {
	Level string
}

func InitLogger(cfg Config) error {
	if level, err := logrus.ParseLevel(cfg.Level); err != nil {
		return err
	} else {
		Logger.SetLevel(level)
	}
	return nil
}
