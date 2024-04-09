package logger

import (
	"github.com/kosalnik/metrics/internal/config"
	"github.com/sirupsen/logrus"
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
