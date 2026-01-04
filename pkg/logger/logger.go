package logger

import (
	"context"
	"log"

	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/pkg/logger/logging"
)

// global logger instance
var logger logging.Logger

func ConfigureFromEnvs(envs *env.Envs) logging.Logger {
	opts := logging.ParseLoggerOptsFromEnvs(envs)
	opts.CallerSkip = 2

	var err error
	logger, err = logging.NewLogger(opts)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	return logger
}

func Configure(level logging.LogLevel, logFormat logging.LogFormat) logging.Logger {
	opts := logging.LoggerOpts{
		Level:      level,
		Format:     logFormat,
		CallerSkip: 2,
	}

	var err error
	logger, err = logging.NewLogger(opts)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	return logger
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func WithCtx(ctx context.Context) logging.Logger {
	return logger.WithCtx(ctx)
}

func WithFields(fields logging.Fields) logging.Logger {
	return logger.WithFields(fields)
}

func Close() error {
	return logger.Close()
}
