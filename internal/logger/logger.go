// Package logger initializes logger.
package logger

import (
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

type Config struct {
	Level string
}

func InitLogger(levelName string) error {
	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		return err
	}
	Logger.SetLevel(level)
	return nil
}
