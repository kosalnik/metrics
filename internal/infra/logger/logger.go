// Package logger initializes logger.
package logger

import (
	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/config"
)

var Logger = logrus.New()

func InitLogger(cfg config.Logger) error {
	if level, err := logrus.ParseLevel(cfg.Level); err != nil {
		return err
	} else {
		Logger.SetLevel(level)
	}
	return nil
}
