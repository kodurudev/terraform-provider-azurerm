package sdk

type Logger interface {
	Info(message string)
	InfoF(format string, args ...interface{})
	Warn(message string)
	WarnF(format string, args ...interface{})
}
