package log

import "context"

// Logger 接口定义记录日志的能力
type Logger interface {
	// Debug 日志
	Debug(msg string, fields ...Field)
	Debugf(format string, v ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	// Info 日志
	Info(msg string, fields ...Field)
	Infof(format string, v ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	// Warn 日志
	Warn(msg string, fields ...Field)
	Warnf(format string, v ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	// Error 日志
	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
	Errorw(msg string, keysAndValues ...interface{})

	// Panic 日志
	Panic(msg string, fields ...Field)
	Panicf(format string, v ...interface{})
	Panicw(msg string, keysAndValues ...interface{})

	// Fatal 日志
	Fatal(msg string, fields ...Field)
	Fatalf(format string, v ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	// WithValues 向 logger 添加一些上下文的 key-value 对。
	WithValues(keysAndValues ...interface{}) Logger

	// WithName 向 logger 的名称添加一个新元素，建议名称段仅包含字母、数字和连字符
	WithName(name string) Logger

	// WithContext 返回设置了日志值的 context 副本
	WithContext(ctx context.Context) context.Context

	// Flush 调用底层 Core 的 Sync 方法，刷新所有缓冲的日志条目，应用程序应注意在退出之前调用 Sync
	Flush()
}
