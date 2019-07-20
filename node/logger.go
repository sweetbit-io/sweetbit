package node

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// A compile time check to ensure that noopLogger fully implements the Logger interface
var _ Logger = (*noopLogger)(nil)

type noopLogger struct {
}

func (l noopLogger) Infof(format string, args ...interface{})  {}
func (l noopLogger) Errorf(format string, args ...interface{}) {}
