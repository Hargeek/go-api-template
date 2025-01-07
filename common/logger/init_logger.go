package logger

import (
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

var (
	logger *slog.Logger
)

func InitLogger() {
	// 设置日志级别
	level := getLogLevel()

	// 设置日志输出
	output := os.Stdout
	if viper.GetString("logging.output") == "stderr" {
		output = os.Stderr
	}

	// 创建日志处理器
	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				// 自定义时间格式
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000"))
			}
			return a
		},
	})

	// 创建日志记录器
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

func BaseLogger() *slog.Logger {
	return logger
}
