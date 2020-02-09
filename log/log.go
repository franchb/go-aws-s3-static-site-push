package log

import (
	"github.com/rs/zerolog"
	"os"
)

var log = zerolog.New(os.Stderr).With().Timestamp().Logger()

// InitLogs initializes zerolog with custom params
func InitLogs(level Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	switch level {
	case LevelInfo: zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LevelWarning:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// SetHumanFriendly to log a human-friendly, colorized output
func SetHumanFriendly() {
	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// Debug Set debug level for logs
func Debug() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

