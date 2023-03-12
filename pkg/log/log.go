package log

import (
	"context"
	"go.uber.org/zap"
	"sync"
)

var (
	std = New(NewOptions())
	mu  sync.Mutex
)

// Init 用指定的 Options 初始化 log 包的全局 logger.
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

// SugaredLogger returns global sugared logger.
func SugaredLogger() *zap.SugaredLogger {
	return std.zapLogger.Sugar()
}

// ZapLogger 用于其它 log 的包装
func ZapLogger() *zap.Logger {
	return std.zapLogger
}

func Debug(msg string, fields ...Field) {
	std.zapLogger.Debug(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Debugf(format, v...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

func Info(msg string, fields ...Field) {
	std.zapLogger.Info(msg, fields...)
}

func Infof(format string, v ...interface{}) {
	std.zapLogger.Sugar().Infof(format, v...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

func Warn(msg string, fields ...Field) {
	std.zapLogger.Warn(msg, fields...)
}

func Warnf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Warnf(format, v...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

func Error(msg string, fields ...Field) {
	std.zapLogger.Error(msg, fields...)
}

func Errorf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Errorf(format, v...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

func Panic(msg string, fields ...Field) {
	std.zapLogger.Panic(msg, fields...)
}

func Panicf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Panicf(format, v...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

func Fatal(msg string, fields ...Field) {
	std.zapLogger.Fatal(msg, fields...)
}

func Fatalf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Fatalf(format, v...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

func WithValues(keysAndValues ...interface{}) Logger {
	return std.WithValues(keysAndValues...)
}

func WithName(s string) Logger {
	return std.WithName(s)
}

func WithContext(ctx context.Context) context.Context {
	return std.WithContext(ctx)
}

func Flush() {
	std.Flush()
}

func L(ctx context.Context) *zapLogger {
	return std.L(ctx)
}
