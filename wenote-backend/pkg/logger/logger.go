package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

// Init 初始化日志
func Init(mode string) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		AddSource: true,
	}

	if mode == "debug" {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		opts.Level = slog.LevelInfo
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

// Info 记录信息日志
func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

// Error 记录错误日志
func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}

// Warn 记录警告日志
func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}
