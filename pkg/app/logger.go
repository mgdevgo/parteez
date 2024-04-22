package app

import (
	"log/slog"
	"os"

	"iditusi/pkg/env"
	"iditusi/pkg/logging/handlers"
)

func newLogger(appEnv string) *slog.Logger {
	logger := new(slog.Logger)

	switch appEnv {
	case env.Local:
		opts := handlers.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		logger = slog.New(opts.NewPrettyHandler(os.Stdout))
	case env.Dev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case env.Prod:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
