package sdk

// NullLogger will debug the log output and is designed only for running the
// debug logger in production to reduce console output
type NullLogger struct {
}

func (NullLogger) Info(_ string) {
}

func (NullLogger) Infof(_ string, _ ...interface{}) {
}

func (NullLogger) Warn(_ string) {
}

func (NullLogger) Warnf(_ string, _ ...interface{}) {
}
