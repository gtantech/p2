package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/gtantech/p2/internal/env"
)

func New(service string, deploymentType env.DeploymentType, file *os.File) *slog.Logger {
	level := slog.LevelInfo

	if deploymentType == env.DEV {
		level = slog.LevelDebug
	}

	w := io.MultiWriter(os.Stdout, file)

	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     level,
		AddSource: deploymentType != env.PROD,
	})

	return slog.New(handler).With(
		"service", service,
		"environment", deploymentType,
	)
}
