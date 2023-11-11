package server

import "github.com/rs/zerolog"

type echoLogger struct {
	level zerolog.Level
	log   zerolog.Logger
}

// Write writes the given bytes to the logger with the given level
func (l *echoLogger) Write(p []byte) (n int, err error) {
	l.log.WithLevel(l.level).Msgf("%s", p)
	return len(p), nil
}
