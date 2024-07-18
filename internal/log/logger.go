// Package logger initializes logger.
package log

import (
	"os"

	"github.com/rs/zerolog"
)

var loggerInstance = zerolog.New(os.Stderr).With().Timestamp().Logger()

type Config struct {
	Level string
}

func InitLogger(levelName string) error {
	level, err := zerolog.ParseLevel(levelName)
	if err != nil {
		return err
	}
	loggerInstance = loggerInstance.Level(level)
	return nil
}

func Debug() *zerolog.Event {
	return loggerInstance.Debug()
}

func Warn() *zerolog.Event {
	return loggerInstance.Warn()
}

func Warning() *zerolog.Event {
	return loggerInstance.Warn()
}

func Info() *zerolog.Event {
	return loggerInstance.Info()
}

func Error() *zerolog.Event {
	return loggerInstance.Error()
}

func Fatal() *zerolog.Event {
	return loggerInstance.Fatal()
}

func Panic() *zerolog.Event {
	return loggerInstance.Panic()
}

func Log() *zerolog.Event {
	return loggerInstance.Log()
}
