package migrator

import "log"

type Logger struct {
	logger log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: *log.Default(),
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *Logger) Verbose() bool {
	return l.logger.Flags()&log.Lshortfile != 0
}
