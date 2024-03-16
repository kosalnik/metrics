package logger

import (
	"github.com/sirupsen/logrus"
)

func InitLogger() error {
	level, err := logrus.ParseLevel("info")
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	return nil
}
