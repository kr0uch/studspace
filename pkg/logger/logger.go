package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	LogFilePath string `yaml:"log_path" env-default:"logs"`
	Level       string `yaml:"level" env-default:"info"`
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger
	Sync() error
}

type ZapLogger struct {
	*zap.Logger
}

func NewLogger(cfg LoggerConfig) (Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(cfg.LogFilePath, 0750)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filepath.Join(cfg.LogFilePath, "server.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return nil, err
	}

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(file),
		level,
	)

	stdout := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		level,
	)

	core := zapcore.NewTee(fileCore, stdout)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))

	return &ZapLogger{Logger: zapLogger}, nil
}

type Field zap.Field

func String(key, val string) Field {
	return Field(zap.String(key, val))
}

func Int(key string, val int) Field {
	return Field(zap.Int(key, val))
}

func Int64(key string, val int64) Field {
	return Field(zap.Int64(key, val))
}

func Bool(key string, val bool) Field {
	return Field(zap.Bool(key, val))
}

func Duration(key string, val time.Duration) Field {
	return Field(zap.Duration(key, val))
}

func Time(key string, val time.Time) Field {
	return Field(zap.Time(key, val))
}

func Any(key string, val any) Field {
	return Field(zap.Any(key, val))
}

func Error(err error) Field {
	return Field(zap.Error(err))
}

func (z *ZapLogger) Debug(msg string, fields ...Field) {
	z.Logger.Debug(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.Logger.Info(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Warn(msg string, fields ...Field) {
	z.Logger.Warn(msg, toZapFields(fields)...)
}

func (z *ZapLogger) Error(msg string, fields ...Field) {
	z.Logger.Error(msg, toZapFields(fields)...)
}

func (z *ZapLogger) With(fields ...Field) Logger {
	return &ZapLogger{Logger: z.Logger.With(toZapFields(fields)...)}
}

func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Field(field)
	}
	return zapFields
}
