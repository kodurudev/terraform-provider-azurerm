package sdk

type Logger interface {
	Info(message string)
	Infof(format string, args ...interface{})
	Warn(message string)
	Warnf(format string, args ...interface{})
}
