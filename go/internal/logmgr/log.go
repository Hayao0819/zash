package logmgr

import (
	"log/slog"

	"github.com/m-mizutani/clog"
)

func defaultHandler() slog.Handler {
	return clog.New(clog.WithColor(true))
}

func debugHandler() slog.Handler {
	return clog.New(clog.WithColor(true), clog.WithLevel(slog.LevelDebug))
}

func SetDefaultLogger() {
	slog.SetDefault(slog.New(defaultHandler()))
}

func SetDebugLogger() {
	slog.SetDefault(slog.New(debugHandler()))
}
