package logger

import "strings"

type MigratorLogger struct{}

func NewMigratorLogger() *MigratorLogger {
	return &MigratorLogger{}
}

func (MigratorLogger) Printf(format string, v ...interface{}) {
	format = strings.TrimSuffix(format, "\n")
	logger.Infof(format, v...)
}

func (MigratorLogger) Verbose() bool {
	return true
}
