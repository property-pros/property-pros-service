package common

import (
	"log"
)

type Logger struct{}

func NewLogger() *Logger {
    return &Logger{}
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
    log.Printf(msg, keysAndValues...)
}