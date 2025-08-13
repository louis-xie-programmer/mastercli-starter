package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(level string) zerolog.Logger {
	lvl := zerolog.InfoLevel
	switch strings.ToLower(level) {
	case "debug":
		lvl = zerolog.DebugLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "error":
		lvl = zerolog.ErrorLevel
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	l := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(lvl)
	log = l
	return l
}

func L() zerolog.Logger {
	return log
}
