package ctxlog

import (
	"log/slog"
)

type Level = slog.Level

const (
	DEBUG = slog.LevelDebug
	INFO  = slog.LevelInfo
	WARN  = slog.LevelWarn
	ERROR = slog.LevelError
)

var (
	Default    = slog.Default
	SetDefault = slog.SetDefault
)
