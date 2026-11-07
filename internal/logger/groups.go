package logger

import "log/slog"

func HttpErrorGroup(status int, errMsg string, args ...any) slog.Attr {
	return slog.Group("http", append([]any{"status", status, "errMsg", errMsg}, args)...)
}
