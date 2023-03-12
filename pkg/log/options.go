package log

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	consoleFormat = "console" // 日志输出格式：console
	jsonFormat    = "json"    // 日志输出格式：json
)

var (
	atomicLevel = zap.NewAtomicLevel()
	changeUrl   = "/change/level"
)

// Options 包含所有日志相关的配置项
type Options struct {
	// Logger 的名字
	Name string `json:"name" mapstructure:"name"`
	// 是否是开发模式。如果是开发模式，会对 WarnLevel 以上进行堆栈跟踪；
	//               非开发模式，会对 PanicLevel 以上进行堆栈跟踪。
	Development bool `json:"development" mapstructure:"development"`
	// 日志级别，优先级从低到高依次为：Debug, Info, Warn, Error, Dpanic, Panic, Fatal
	Level string `json:"level" mapstructure:"level"`
	// 是否开启颜色输出，true，是；false，否
	EnableColor bool `json:"enable-color" mapstructure:"enable-color"`
	// 是否开启 caller，如果开启会在日志中显示调用日志所在的文件、函数和行号
	DisableCaller bool `json:"disable-caller" mapstructure:"disable-caller"`
	// 是否在 WarnLevel（开发模式） 或 PanicLevel（非开发模式） 及以上级别禁止打印堆栈信息
	DisableStacktrace bool `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	// 支持的日志输出格式，支持 Console 和 JSON 两种
	Format string `json:"format" mapstructure:"format"`
	// 支持输出到多个输出，用逗号分开。支持输出到标准输出（stdout）和文件
	OutputPaths []string `json:"output-paths" mapstructure:"output-paths"`
	// zap 内部 (非业务) 错误日志输出路径，多个输出，用逗号分开
	ErrorOutputPaths []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	// MaxSize 是日志文件轮转之前的最大大小（MB）。默认为 100MB。
	MaxSize int `json:"maxsize" mapstructure:"maxsize"`
	// MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数。
	MaxAge int `json:"maxage" mapstructure:"maxage"`
	// MaxBackups 是要保留的旧日志文件的最大数量。
	MaxBackups int `json:"maxbackups" mapstructure:"maxbackups"`
	// Compress 确定是否应使用 gzip 压缩轮转的日志文件。
	Compress bool `json:"compress" mapstructure:"compress"`
	// Port 用于指定动态修改日志输出级别端口，默认为 0 （不开启）。
	Port int `json:"port" mapstructure:"port"`
}

func NewOptions() *Options {
	return &Options{
		Development:       false,
		Level:             zapcore.InfoLevel.String(),
		EnableColor:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            consoleFormat,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		MaxSize:           10,
		MaxAge:            7,
		MaxBackups:        5,
		Compress:          true,
	}
}

// Validate 验证 options 字段
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Build 从 Config 和 Options 构建 zap 包的 logger。
func (o *Options) Build() error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		// 解析出错， 默认为 info 级别
		zapLevel = zapcore.InfoLevel
	}

	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == consoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stacktrace",
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	enc := zapcore.NewConsoleEncoder(encoderConfig)
	if o.Format == jsonFormat {
		enc = zapcore.NewJSONEncoder(encoderConfig)
	}

	writer := getWriter(o, o.OutputPaths...)
	errWriter := getWriter(o, o.ErrorOutputPaths...)
	atomicLevel.SetLevel(zapLevel)
	l := zap.New(zapcore.NewCore(enc, writer, atomicLevel),
		o.buildOptions(errWriter)...)
	l = l.Named(o.Name)

	zap.RedirectStdLog(l)
	zap.ReplaceGlobals(l)

	if o.Port != 0 {
		http.HandleFunc(changeUrl, atomicLevel.ServeHTTP)
		go func() {
			addr := fmt.Sprintf(":%d", o.Port)
			if err := http.ListenAndServe(addr, nil); err != nil {
				panic(err)
			}
		}()
	}

	return nil
}

func (o *Options) buildOptions(errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}

	if o.Development {
		opts = append(opts, zap.Development())
	}

	if !o.DisableCaller {
		opts = append(opts, zap.AddCaller())
		opts = append(opts, zap.AddCallerSkip(1))
	}

	stackLevel := PanicLevel
	if o.Development {
		stackLevel = WarnLevel
	}
	if !o.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	return opts
}

// 时间格式
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

func getWriter(opts *Options, paths ...string) zapcore.WriteSyncer {
	writers := make([]zapcore.WriteSyncer, 0, len(paths))

	for _, path := range paths {
		if path == "stdout" {
			writers = append(writers, zapcore.AddSync(os.Stdout))
		} else if path == "stderr" {
			writers = append(writers, zapcore.AddSync(os.Stderr))
		} else {
			lumberjackLogger := &lumberjack.Logger{
				Filename:   path,
				MaxSize:    opts.MaxSize,
				MaxAge:     opts.MaxAge,
				MaxBackups: opts.MaxBackups,
				LocalTime:  true,
				Compress:   opts.Compress,
			}
			writers = append(writers, zapcore.AddSync(lumberjackLogger))
		}
	}
	writer := zap.CombineWriteSyncers(writers...)

	return writer
}
