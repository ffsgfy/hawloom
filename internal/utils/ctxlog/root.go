package ctxlog

import (
	"context"
	"io"
	"log/slog"
)

func New(writer io.Writer, level Level) *slog.Logger {
	return slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level}))
}

type loggerKeyType struct{}

var loggerKey loggerKeyType

func Logger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return Default()
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func With(ctx context.Context, args ...any) context.Context {
	return WithLogger(ctx, Logger(ctx).With(args...))
}

func WithGroup(ctx context.Context, name string) context.Context {
	return WithLogger(ctx, Logger(ctx).WithGroup(name))
}

func Log(ctx context.Context, level Level, msg string, args ...any) {
	Logger(ctx).Log(ctx, level, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	Logger(ctx).DebugContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	Logger(ctx).InfoContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	Logger(ctx).WarnContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	Logger(ctx).ErrorContext(ctx, msg, args...)
}

func Error2(ctx context.Context, msg string, err error, args ...any) {
	Error(ctx, msg, append(args, "err", err)...)
}
