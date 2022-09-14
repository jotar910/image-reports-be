package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Initialize() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile("tmp/image-reports.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.FatalLevel))
}

// Implement zap logger

func Named(s string) *zap.Logger {
	return logger.Named(s)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return logger.WithOptions(opts...)
}

func With(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

func Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return logger.Check(lvl, msg)
}

func Log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	logger.Log(lvl, msg, fields...)
}

func Logf(lvl zapcore.Level, msg string, a ...any) {
	logger.Log(lvl, fmt.Sprintf(msg, a...))
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Debugf(msg string, a ...any) {
	logger.Debug(fmt.Sprintf(msg, a...))
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Infof(msg string, a ...any) {
	logger.Info(fmt.Sprintf(msg, a...))
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Warnf(msg string, a ...any) {
	logger.Warn(fmt.Sprintf(msg, a...))
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Errorf(msg string, a ...any) {
	logger.Error(fmt.Sprintf(msg, a...))
}

func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

func DPanicf(msg string, a ...any) {
	logger.DPanic(fmt.Sprintf(msg, a...))
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Panicf(msg string, a ...any) {
	logger.Panic(fmt.Sprintf(msg, a...))
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Fatalf(msg string, a ...any) {
	logger.Fatal(fmt.Sprintf(msg, a...))
}
