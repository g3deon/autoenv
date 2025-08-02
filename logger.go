package autoenv

import "fmt"

type Logger interface {
	InfoF(format string, args ...any)
	WarnF(format string, args ...any)
	DebugF(format string, args ...any)
	ErrorF(format string, args ...any)
}

type defaultLogger struct{}

func (l *defaultLogger) InfoF(format string, args ...any) {
	fmt.Println("[AUTO ENV] INFO: ", fmt.Sprintf(format, args...))
}

func (l *defaultLogger) WarnF(format string, args ...any) {
	fmt.Println("[AUTO ENV] WARN: ", fmt.Sprintf(format, args...))
}

func (l *defaultLogger) DebugF(format string, args ...any) {
	fmt.Println("[AUTO ENV] DEBUG: ", fmt.Sprintf(format, args...))
}

func (l *defaultLogger) ErrorF(format string, args ...any) {
	fmt.Println("[AUTO ENV] ERROR: ", fmt.Sprintf(format, args...))
}
