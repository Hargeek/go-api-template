package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"go-api-template/common/types"

	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/trace"
)

var (
	// 默认指向 slog.Default()，保证 InitLogger 调用前（如单测场景）也不会 nil panic
	logger      = slog.Default()
	fileWriters []*os.File // 保存打开的文件句柄，避免被GC回收
)

func InitLogger() {
	// 关闭上一次初始化打开的文件句柄（支持多次调用）
	for _, f := range fileWriters {
		_ = f.Close()
	}
	fileWriters = nil

	level := getLogLevel()

	// 从配置读取输出列表
	outputs := viper.GetStringSlice("logging.output")
	if len(outputs) == 0 {
		// 如果没有配置，默认使用stdout
		outputs = []string{"stdout"}
	}

	// 构建多个输出目标
	writers := make([]io.Writer, 0, len(outputs))
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

	localHandler := slog.NewJSONHandler(outputWriter, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000"))
			}
			return a
		},
	})

	handlers := []slog.Handler{localHandler}

	// OTEL_EXPORTER_ENABLED_LOGS=true 时追加 otelslog bridge handler
	// 未配置 OTEL_EXPORTER_OTLP_ENDPOINT 时 LoggerProvider 为 noop，handler 不产生额外输出
	// telemetry.Setup() 必须在 InitLogger() 之前调用，确保 LoggerProvider 已注册
	if os.Getenv("OTEL_EXPORTER_ENABLED") == "true" &&
		os.Getenv("OTEL_EXPORTER_ENABLED_LOGS") == "true" {
		handlers = append(handlers, otelslog.NewHandler(types.ServiceName))
	}

	logger = slog.New(newMultiHandler(handlers...))
	slog.SetDefault(logger)
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

// withTrace 从 context 中提取 otel span，将 trace_id / span_id 追加到日志参数中
// 作为本地 JSON 输出的 fallback：无论 OTEL Log 是否启用，本地日志始终携带 trace_id
// otelslog bridge 启用时，OTEL 侧会从 context 自动关联，两者互不干扰
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

// multiHandler 将日志记录分发给多个 slog.Handler
type multiHandler struct {
	handlers []slog.Handler
}

func newMultiHandler(handlers ...slog.Handler) slog.Handler {
	if len(handlers) == 1 {
		return handlers[0]
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
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
	panic(msg)
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
