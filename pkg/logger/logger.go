package logger

import (
	"log/slog"
	"os"

	"parteez/pkg/env"
	"parteez/pkg/logger/handlers"
)

func New(appenv string) *slog.Logger {
	log := new(slog.Logger)

	switch appenv {
	case env.Local:
		opts := handlers.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		log = slog.New(opts.NewPrettyHandler(os.Stdout))
	case env.Dev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case env.Prod:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
