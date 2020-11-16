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

func (ConsoleLogger) InfoF(format string, args ...interface{}) {
	log.Print(fmt.Sprintf(format, args...))
}

func (ConsoleLogger) Warn(message string) {
	log.Print(message)
}

func (ConsoleLogger) WarnF(format string, args ...interface{}) {
	log.Print(fmt.Sprintf(format, args...))
}
