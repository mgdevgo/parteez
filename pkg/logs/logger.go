package logs

import (
	"log/slog"
	"os"

	"iditusi/pkg/core"
	"iditusi/pkg/logs/handlers"
)

type Ctx struct {
}

func NewLogger(env core.AppEnv) *slog.Logger {
	logger := new(slog.Logger)

	switch env {
	case core.AppEnvLocal:
		opts := handlers.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		logger = slog.New(opts.NewPrettyHandler(os.Stdout))
	case core.AppEnvDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case core.AppEnvProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
