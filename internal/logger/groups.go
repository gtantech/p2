package logger

import "log/slog"

func HttpErrorGroup(status int, errMsg string, args ...any) slog.Attr {
	groupName := "http"
	attr := []any{"status", status, "errMsg", errMsg}
	if len(args) == 0 {
		return slog.Group(groupName, attr...)
	}
	return slog.Group(groupName, append(attr, args)...)
}

func JsonGroup(json string) slog.Attr {
	return slog.Group("json", "value", json)
}
