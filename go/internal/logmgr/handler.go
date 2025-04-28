package logmgr

import (
	"log/slog"

	"github.com/m-mizutani/clog"
)

func colorHandler() slog.Handler {
	return clog.New(clog.WithColor(true), clog.WithLevel(slog.LevelDebug))
}

func nopHandler() slog.Handler {
	return clog.New(clog.WithLevel(slog.Level(100)))
}

func categorizedColorLogger(c string) *slog.Logger {
	return slog.New(colorHandler()).With("category", c)
}

var shellLogger *slog.Logger = categorizedColorLogger("shell")
var lexerLogger *slog.Logger = categorizedColorLogger("lexer")
var parserLogger *slog.Logger = categorizedColorLogger("parser")
var builtinLogger *slog.Logger = categorizedColorLogger("builtin")
var logmgrLogger *slog.Logger = categorizedColorLogger("logmgr")
var cmdLogger *slog.Logger = categorizedColorLogger("cmd")
var cmdExecLogger *slog.Logger = categorizedColorLogger("cmdexec")
var nopLogger *slog.Logger = slog.New(nopHandler())

var showLogs bool = false

func EnableLogs() {
	showLogs = true
}
func DisableLogs() {
	showLogs = false
}

func Shell() *slog.Logger {
	if showLogs {
		return shellLogger
	}
	return nopLogger
}
func Lexer() *slog.Logger {
	if showLogs {
		return lexerLogger
	}
	return nopLogger
}
func Parser() *slog.Logger {
	if showLogs {
		return parserLogger
	}
	return nopLogger
}
func Builtin() *slog.Logger {
	if showLogs {
		return builtinLogger
	}
	return nopLogger
}
func Logmgr() *slog.Logger {
	if showLogs {
		return logmgrLogger
	}
	return nopLogger
}
func Cmd() *slog.Logger {
	if showLogs {
		return cmdLogger
	}
	return nopLogger
}
func CmdExec() *slog.Logger {
	if showLogs {
		return cmdExecLogger
	}
	return nopLogger
}
