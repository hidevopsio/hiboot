package jaeger

import "hidevops.io/hiboot/pkg/log"

// StdLogger is implementation of the Logger interface that delegates to default `log` package
type Logger struct{}

func (l *Logger) Error(msg string) {
	log.Errorf("ERROR: %s", msg)
}

// Infof logs a message at info priority
func (l *Logger) Infof(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}
