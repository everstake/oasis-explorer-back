package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

const (
	LevelDebug = "debug"
	LevelWarn  = "warn"
	LevelInfo  = "info"
	LevelError = "error"
)

func getZapLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		panic(fmt.Sprintf("wrong log level %s ", level))
	}
}

func init() {
	logger = getLogger(zapcore.DebugLevel)
}

func SetLogLevel(logLevel string) {
	logger = getLogger(getZapLevel(logLevel))
}

func getLogger(logLevel zapcore.Level) *zap.Logger {
	var err error
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(logLevel)
	cfg.Encoding = "console"
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	logger, err = cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("Can`t init logger: %s", err.Error()))
	}
	return logger
}

func Info(text string, fields ...zap.Field) {
	logger.Info(text, fields...)
}

func Error(text string, fields ...zap.Field) {
	logger.Error(text, fields...)
}

func Warn(text string, fields ...zap.Field) {
	logger.Warn(text, fields...)
}

func Debug(text string, fields ...zap.Field) {
	logger.Debug(text, fields...)
}

func Fatal(text string, fields ...zap.Field) {
	logger.Fatal(text, fields...)
}
