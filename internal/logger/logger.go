package logger

import (
	"log/slog"
	"os"

	"github.com/gtantech/p2/internal/env"
)

func New(service string, deploymentType env.DeploymentType) *slog.Logger {
	level := slog.LevelInfo

	if deploymentType == env.DEV {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: deploymentType != env.PROD,
	})

	return slog.New(handler).With(
		"service", service,
		"environment", deploymentType,
	)
}
