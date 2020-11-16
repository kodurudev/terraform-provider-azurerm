package sdk

import (
	"fmt"
	"log"
)

type ConsoleLogger struct {
}

func (ConsoleLogger) Info(message string) {
	log.Print(message)
}

func (ConsoleLogger) Infof(format string, args ...interface{}) {
	log.Print(fmt.Sprintf(format, args...))
}

func (ConsoleLogger) Warn(message string) {
	log.Print(message)
}

func (ConsoleLogger) Warnf(format string, args ...interface{}) {
	log.Print(fmt.Sprintf(format, args...))
}
