package logging

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/amahdian/cliplab-be/global/env"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

var logLevels = []LogLevel{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel}

type LogFormat string

const (
	TextFormat LogFormat = "text"
	JsonFormat LogFormat = "json"
)

var logFormats = []LogFormat{TextFormat, JsonFormat}

// Logger is the interface that wraps the basic logging methods.
// The goal is to provide a simple and consistent logging interface
// that can be implemented by different logging libraries.
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	WithCtx(ctx context.Context) Logger
	WithFields(map[string]any) Logger

	IsLevelEnabled(level LogLevel) bool

	Close() error
}

type zapLogger struct {
	ctx    context.Context
	logger *otelzap.Logger
	fields Fields
}

type LoggerOpts struct {
	Level      LogLevel
	Format     LogFormat
	CallerSkip int
}

func ParseLoggerOptsFromEnvs(envs *env.Envs) LoggerOpts {
	logLevel := LogLevel(strings.ToLower(envs.Server.LogLevel))
	logFormat := LogFormat(strings.ToLower(envs.Server.LogFormat))

	return LoggerOpts{
		Level:      logLevel,
		Format:     logFormat,
		CallerSkip: 1,
	}
}

func NewLogger(opts LoggerOpts) (Logger, error) {
	var logLevel zapcore.Level
	switch opts.Level {
	case DebugLevel:
		logLevel = zapcore.DebugLevel
	case InfoLevel:
		logLevel = zapcore.InfoLevel
	case WarnLevel:
		logLevel = zapcore.WarnLevel
	case ErrorLevel:
		logLevel = zapcore.ErrorLevel
	case FatalLevel:
		logLevel = zapcore.FatalLevel
	default:
		return nil, fmt.Errorf("Log level is not one of the supported values (%s): %s", logLevels, logLevel)
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = timeEncoder
	var encoder zapcore.Encoder
	switch opts.Format {
	case "", "text":
		encoder = zapcore.NewConsoleEncoder(config)
	case "json":
		encoder = zapcore.NewJSONEncoder(config)
	default:
		return nil, fmt.Errorf("Log format is not one of the supported values (%s): %s", logFormats, opts.Format)
	}

	// TODO: allow multiple outputs like file
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel),
	)
	zLogger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(opts.CallerSkip),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	otelLogger := otelzap.New(zLogger, otelzap.WithMinLevel(logLevel))

	return &zapLogger{
		logger: otelLogger,
		fields: make(map[string]any),
	}, nil
}

func (z *zapLogger) Debug(args ...interface{}) {
	if z.ctx != nil {
		z.logger.DebugContext(z.ctx, fmt.Sprint(args...), z.fields.asZapFields()...)
	}
	z.logger.Debug(fmt.Sprint(args...), z.fields.asZapFields()...)
}

func (z *zapLogger) Info(args ...interface{}) {
	if z.ctx != nil {
		z.logger.InfoContext(z.ctx, fmt.Sprint(args...), z.fields.asZapFields()...)
	}
	z.logger.Info(fmt.Sprint(args...), z.fields.asZapFields()...)
}

func (z *zapLogger) Warn(args ...interface{}) {
	if z.ctx != nil {
		z.logger.WarnContext(z.ctx, fmt.Sprint(args...), z.fields.asZapFields()...)
	}
	z.logger.Warn(fmt.Sprint(args...), z.fields.asZapFields()...)
}

func (z *zapLogger) Error(args ...interface{}) {
	if z.ctx != nil {
		z.logger.ErrorContext(z.ctx, fmt.Sprint(args...), z.fields.asZapFields()...)
	}
	z.logger.Error(fmt.Sprint(args...), z.fields.asZapFields()...)
}

func (z *zapLogger) Fatal(args ...interface{}) {
	if z.ctx != nil {
		z.logger.FatalContext(z.ctx, fmt.Sprint(args...), z.fields.asZapFields()...)
	}
	z.logger.Fatal(fmt.Sprint(args...), z.fields.asZapFields()...)
}

func (z *zapLogger) Debugf(format string, args ...interface{}) {
	if z.ctx != nil {
		z.logger.DebugContext(z.ctx, fmt.Sprintf(format, args...), z.fields.asZapFields()...)
	}
	z.logger.Debug(fmt.Sprintf(format, args...), z.fields.asZapFields()...)
}

func (z *zapLogger) Infof(format string, args ...interface{}) {
	if z.ctx != nil {
		z.logger.InfoContext(z.ctx, fmt.Sprintf(format, args...), z.fields.asZapFields()...)
	}
	z.logger.Info(fmt.Sprintf(format, args...), z.fields.asZapFields()...)
}

func (z *zapLogger) Warnf(format string, args ...interface{}) {
	if z.ctx != nil {
		z.logger.WarnContext(z.ctx, fmt.Sprintf(format, args...), z.fields.asZapFields()...)
	}
	z.logger.Warn(fmt.Sprintf(format, args...), z.fields.asZapFields()...)
}

func (z *zapLogger) Errorf(format string, args ...interface{}) {
	if z.ctx != nil {
		z.logger.ErrorContext(z.ctx, fmt.Sprintf(format, args...), z.fields.asZapFields()...)
	}
	z.logger.Error(fmt.Sprintf(format, args...), z.fields.asZapFields()...)
}

func (z *zapLogger) Fatalf(format string, args ...interface{}) {
	if z.ctx != nil {
		z.logger.FatalContext(z.ctx, fmt.Sprintf(format, args...), z.fields.asZapFields()...)
	}
	z.logger.Fatal(fmt.Sprintf(format, args...), z.fields.asZapFields()...)
}

func (z *zapLogger) IsLevelEnabled(level LogLevel) bool {
	var logLevel zapcore.Level
	switch level {
	case DebugLevel:
		logLevel = zapcore.DebugLevel
	case InfoLevel:
		logLevel = zapcore.InfoLevel
	case WarnLevel:
		logLevel = zapcore.WarnLevel
	case ErrorLevel:
		logLevel = zapcore.ErrorLevel
	case FatalLevel:
		logLevel = zapcore.FatalLevel
	default:
		return false
	}
	return z.logger.Core().Enabled(logLevel)
}

func (z *zapLogger) WithCtx(ctx context.Context) Logger {
	return &zapLogger{
		ctx:    ctx,
		logger: z.logger,
		fields: z.fields,
	}
}

func (z *zapLogger) WithFields(fields map[string]any) Logger {
	mergedFields := Fields{}
	mergedFields.mergeIn(fields)
	mergedFields.mergeIn(z.fields)
	return &zapLogger{
		logger: z.logger,
		fields: mergedFields,
	}
}

func (z *zapLogger) Close() error {
	return z.logger.Sync()
}

type Fields map[string]any

func (f Fields) mergeIn(fs ...map[string]any) {
	for _, other := range fs {
		for k, v := range other {
			f[k] = v
		}
	}
}

func (f Fields) asMap() map[string]any {
	return f
}

func (f Fields) asZapFields() []zap.Field {
	zapFields := make([]zap.Field, 0, len(f))
	for k, v := range f {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.Local()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	micros := t.Nanosecond() / 1000

	buf := make([]byte, 27)

	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = 'T'
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((micros/100000)%10) + '0'
	buf[21] = byte((micros/10000)%10) + '0'
	buf[22] = byte((micros/1000)%10) + '0'
	buf[23] = byte((micros/100)%10) + '0'
	buf[24] = byte((micros/10)%10) + '0'
	buf[25] = byte((micros)%10) + '0'
	buf[26] = 'Z'

	enc.AppendString(string(buf))
}
