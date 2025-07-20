package autoenv

import "fmt"

type Logger interface {
	Debugf(format string, args ...any)
	Errorf(format string, args ...any)
}

type defaultLogger struct{}

func (l *defaultLogger) Debugf(format string, args ...any) {
	fmt.Printf("[AUTO ENV] "+format, args...)
}

func (l *defaultLogger) Errorf(format string, args ...any) {
	fmt.Printf("[AUTO ENV] ERROR: "+format, args...)
}
