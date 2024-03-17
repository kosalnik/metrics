package logger

import (
	"github.com/kosalnik/metrics/internal/config"
	"github.com/sirupsen/logrus"
)

func InitLogger(cfg config.Logger) error {
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	return nil
}
