package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/ipcdev/component/pkg/log"
)

func main() {
	// log 包开箱即用
	//log.Debug("this is a debug message") // 默认日志级别为 Info，不会输出 Debug 级别日志
	//
	//log.Info("this is a info message")
	//log.Infof("this is a %s message", "infof")
	//log.Infow("this is a infow message", "username", "admin")
	//
	//log.Warn("this is a warn message")

	fmt.Printf("\n ================== user defined options ===============\n\n")

	opts := &log.Options{
		//Name:              "ipc-logger",
		Development:       false,                           // 是否为开发模式
		Level:             "debug",                         // debug、info、warn、error、painc、fatal
		EnableColor:       false,                           // 是否开启颜色
		DisableCaller:     false,                           // 是否关闭行号、调用者信息
		DisableStacktrace: false,                           // 禁止打印堆栈
		Format:            "console",                       // 支持 “console”、“json” 两种方式
		OutputPaths:       []string{"test.log", "stdout"},  // 输出到 test.log 和 标准输出
		ErrorOutputPaths:  []string{"error.log", "stderr"}, // 输出到 error.log 和 标准错误
		MaxBackups:        5,                               // 最大保留 5 个日志文件备份
		MaxAge:            1,                               // 日志文件保留 1d
		MaxSize:           10,                              // 日志文件最大 10M
		Compress:          true,                            // 开启日志压缩
		// curl -XPUT -d "level=info" http://localhost:9090/change/level （设置日志级别）
		// curl http://localhost:9090/change/level （查看当前日志级别）
		Port: 9090,
	}
	errors := opts.Validate()
	if len(errors) != 0 {
		panic(errors)
	}

	// 使用自定义 options 重新初始化 log 包
	log.Init(opts)
	defer log.Flush()

	// Debug
	log.Debug("This is a user defined Debug message")
	log.Debugf("This is a user defined %s message", "Debugf")
	log.Debugw("This is a user defined Debugw message", "username", "admin")

	// Info
	log.Info("This is a user defined Info message")
	log.Infof("This is a user defined %s message", "Infof")
	log.Infow("This is a user defined Infow message", log.Int32("int_key", 50))

	// Warn
	log.Warn("This is a user defined Warn message")
	log.Warnf("This is a user defined %s message", "Warnf")
	log.Warnw("This is a user defined Warnw message", log.Int32("int_key", 50))

	// Error
	log.Error("This is a user defined Error message")
	log.Errorf("This is a user defined %s message", "Errorf")
	log.Errorw("This is a user defined Errorw message", log.Int32("int_key", 50))

	// Panic
	//log.Panic("This is a panic message")
	//log.Panicf("This is a %s message", "Panicf")
	//log.Panicw("This is a panicw message", "logger-name", opts.Name)

	// Fatal
	//log.Fatal("This is a fatal message")
	//log.Fatalf("This is a %s message", "Fatalf")
	//log.Fatalw("This is a Fatalw message", "logger-name", opts.Name)

	// WithName
	ln := log.WithName("withName-test")
	ln.Info("[WithName] logger")

	// WithValues
	lv := log.WithValues("username", "admin")
	lv.Info("[WithValues] logger")

	lvv := lv.WithValues("Request-ID", "123456")
	lvv.Warn("[WithValues] logger")

	// WithContext
	ctx := log.WithContext(context.Background())
	lc := log.FromContext(ctx)
	lc.Info("[WithContext] logger")

	lNo := log.FromContext(context.Background())
	lNo.Info("No Message printed with [WithContext] logger")

	log.Warn("Warn logger")

	// L()
	ctx1 := context.Background()
	ctx1 = context.WithValue(ctx1, log.KeyRequestID, "123456")
	ctx1 = context.WithValue(ctx1, log.KeyUsername, "admin")
	log.L(ctx1).Info("This is a L logger")

	// 测试日志文件滚动
	go func() {
		for i := 0; i < math.MaxInt; i++ {
			//if i%1000 == 0 {
			//	time.Sleep(time.Second)
			//}
			time.Sleep(time.Second * 10)

			log.Debug("This is a user defined Debug message", log.Int("int_key", i))
			log.Debugf("This is a user defined %s message", "Debugf", log.Int("int_key", i))
			log.Debugw("This is a user defined Debugw message", "username", "admin", log.Int("int_key", i))

			// Info
			log.Info("This is a user defined Info message", log.Int("int_key", i))
			log.Infof("This is a user defined %s message", "Infof", log.Int("int_key", i))
			log.Infow("This is a user defined Infow message", log.Int("int_key", i))

			// Warn
			log.Warn("This is a user defined Warn message", log.Int("int_key", i))
			log.Warnf("This is a user defined %s message", "Warnf", log.Int("int_key", i))
			log.Warnw("This is a user defined Warnw message", log.Int("int_key", i))

			log.Error("This is a user defined Error message", log.Int("int_key", i))
			log.Errorf("This is a user defined %s message", "Errorf", log.Int("int_key", i))
			log.Errorw("This is a user defined Errorw message", log.Int("int_key", i))
		}
	}()
	// 主线程如果提前退出，可能 lumberjack 正在压缩，会导致文件压缩失败。所以设置主线程阻塞
	time.Sleep(time.Second * 3600)
}
