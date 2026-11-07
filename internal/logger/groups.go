package logger

import "log/slog"

func HttpErrorGroup(status int, args ...any) slog.Attr {
	return slog.Group("http", append([]any{"status", status}, args)...)
}
