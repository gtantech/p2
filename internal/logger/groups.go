package logger

import "log/slog"

func HttpErrorGroup(status int, errMsg string) slog.Attr {
	return slog.Group("http", "status", status, "err", errMsg)
}
