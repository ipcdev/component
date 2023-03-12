package log

import (
	"context"
	"go.uber.org/zap"
)

var _ Logger = (*zapLogger)(nil)

type zapLogger struct {
	zapLogger *zap.Logger
}

// New 通过 opts 新建 logger。
func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}
	opts.Build()

	logger := &zapLogger{
		zapLogger: zap.L(),
	}

	return logger
}

// NewLogger 使用给定的 Zap Logger 来创建一个 Logger
func NewLogger(l *zap.Logger) Logger {
	return &zapLogger{
		zapLogger: l,
	}
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *zapLogger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.zapLogger.Panic(msg, fields...)
}

func (l *zapLogger) Panicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Panicf(format, v...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.zapLogger.Fatal(msg, fields...)
}

func (l *zapLogger) Fatalf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Fatalf(format, v...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) WithValues(keysAndValues ...interface{}) Logger {
	newLogger := l.zapLogger.With(handleFields(l.zapLogger, keysAndValues)...)

	return NewLogger(newLogger)
}

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.zapLogger.Named(name)

	return NewLogger(newLogger)
}

func (l *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, l)

}

func (l *zapLogger) Flush() {
	_ = l.zapLogger.Sync()
}

// L 可以很方便地从 Context 中提取出指定的 key-value，作为上下文添加到日志输出中。
// 当前支持的为 KeyRequestID、KeyUsername
func (l *zapLogger) L(ctx context.Context) *zapLogger {
	lg := l.clone()

	if requestID := ctx.Value(KeyRequestID); requestID != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyRequestID, requestID))
	}
	if username := ctx.Value(KeyUsername); username != nil {
		lg.zapLogger = lg.zapLogger.With(zap.Any(KeyUsername, username))
	}

	return lg
}

func (l *zapLogger) clone() *zapLogger {
	copy := *l

	return &copy
}

// handleFields 将一组任意 key-value 对转换为 Zap字段
func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	if len(args) == 0 {
		return additional
	}

	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))
			break
		}

		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))

			break
		}

		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)

			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}
