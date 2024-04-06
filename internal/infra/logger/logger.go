package logger

import (
	"github.com/sirupsen/logrus"

	"github.com/kosalnik/metrics/internal/config"
)

func InitLogger(cfg config.Logger) error {
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	return nil
}
