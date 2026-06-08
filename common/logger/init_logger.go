package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
)

var (
	logger      *slog.Logger
	fileWriters []*os.File // 保存打开的文件句柄，避免被GC回收
)

func InitLogger() {
	level := getLogLevel()

	// 从配置读取输出列表
	outputs := viper.GetStringSlice("logging.output")
	if len(outputs) == 0 {
		// 如果没有配置，默认使用stdout
		outputs = []string{"stdout"}
	}

	// 构建多个输出目标
	writers := make([]io.Writer, 0)
	for _, output := range outputs {
		switch output {
		case "stdout":
			writers = append(writers, os.Stdout)
		case "stderr":
			writers = append(writers, os.Stderr)
		default:
			// 文件输出
			file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				panic(fmt.Sprintf("Failed to open log file %s: %v", output, err))
			}
			fileWriters = append(fileWriters, file)
			writers = append(writers, file)
		}
	}

	// 使用MultiWriter合并多个输出
	var outputWriter io.Writer = os.Stdout // 默认值
	if len(writers) > 0 {
		outputWriter = io.MultiWriter(writers...)
	}

	handler := slog.NewJSONHandler(outputWriter, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000"))
			}
			return a
		},
	})

	logger = slog.New(handler)
}

func getLogLevel() slog.Level {
	switch viper.GetString("logging.level") {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// withTrace 从 context 中提取 otel span，将 trace_id / span_id 追加到日志参数中。
// Trace 未启用或 span 无效时不追加任何字段。
func withTrace(ctx context.Context, args []interface{}) []interface{} {
	sc := trace.SpanFromContext(ctx).SpanContext()
	if !sc.IsValid() {
		return args
	}
	return append(args,
		slog.String("trace_id", sc.TraceID().String()),
		slog.String("span_id", sc.SpanID().String()),
	)
}

// --- 不带 context 的函数（用于启动/关闭等无请求上下文的场景）---

func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	logger.Debug(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logger.Error(msg, args...)
	os.Exit(1)
}

func Panic(msg string, args ...interface{}) {
	logger.Error(msg, args...)
	panic(fmt.Sprintf(msg, args...))
}

// --- 带 context 的函数（用于请求链路中，自动附加 trace_id / span_id）---

func InfoContext(ctx context.Context, msg string, args ...interface{}) {
	logger.InfoContext(ctx, msg, withTrace(ctx, args)...)
}

func WarnContext(ctx context.Context, msg string, args ...interface{}) {
	logger.WarnContext(ctx, msg, withTrace(ctx, args)...)
}

func ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	logger.ErrorContext(ctx, msg, withTrace(ctx, args)...)
}

func DebugContext(ctx context.Context, msg string, args ...interface{}) {
	logger.DebugContext(ctx, msg, withTrace(ctx, args)...)
}

func BaseLogger() *slog.Logger {
	return logger
}
