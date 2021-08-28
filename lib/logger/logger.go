package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = New()

func New() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.CallerKey = "line"
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg := zap.NewProductionConfig()
	cfg.Encoding = setEncoding("")
	cfg.Level = SetLevel("debug")
	cfg.DisableStacktrace = true
	cfg.EncoderConfig = encoderConfig
	cfg.InitialFields = initFields()
	cfg.OutputPaths = setOutput()
	logger, err := cfg.Build(
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic(err)
	}
	return logger
}

func setOutput(path ...string) []string {
	if len(path) == 0 {
		return []string{"stderr"}
	}
	return path
}

func setEncoding(encoding string) string {
	if encoding == "" {
		return "console"
	}
	return encoding
}

func initFields() map[string]interface{} {
	return map[string]interface{}{}
}

func Debug(msg string, field ...zap.Field) {
	logger.Debug(msg, field...)
}

func Warn(msg string, field ...zap.Field) {
	logger.Warn(msg, field...)
}

func Info(msg string, field ...zap.Field) {
	logger.Info(msg, field...)
}

func Error(msg string, field ...zap.Field) {
	logger.Error(msg, field...)
}

func Errorf(format string, f ...interface{}) {
	logger.Error(fmt.Sprintf(format, f...))
}

func Warnf(format string, f ...interface{}) {
	logger.Warn(fmt.Sprintf(format, f...))
}

func Infof(format string, f ...interface{}) {
	logger.Info(fmt.Sprintf(format, f...))
}

func Debugf(format string, f ...interface{}) {
	logger.Debug(fmt.Sprintf(format, f...))
}

func DebugAsJson(value interface{}) {
	logger.Debug("debugAsJson", zap.Any("object", value))
}

func Err(err error) zap.Field {
	return zap.Error(err)
}

func String(key string, val string) zap.Field {
	return zap.String(key, val)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func Binary(key string, val []byte) zap.Field {
	return zap.Binary(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func ByteString(key string, val []byte) zap.Field {
	return zap.ByteString(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

func Float32(key string, val float32) zap.Field {
	return zap.Float32(key, val)
}

func Int(key string, val int) zap.Field {
	return Int64(key, int64(val))
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Int8(key string, val int8) zap.Field {
	return zap.Int8(key, val)
}

func Uint(key string, val uint) zap.Field {
	return Uint64(key, uint64(val))
}

func Uint64(key string, val uint64) zap.Field {
	return zap.Uint64(key, val)
}

func Uint8(key string, val uint8) zap.Field {
	return zap.Uint8(key, val)
}

func SetLevel(level string) zap.AtomicLevel {
	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		return zap.NewAtomicLevelAt(zap.PanicLevel)
	default:
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	}
}
