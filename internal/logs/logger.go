package logs

import (
	"log/slog"
	"os"

	"iditusi/internal/models"
	"iditusi/pkg/logs/handlers"
)

type Ctx struct {
}

func NewLogger(env models.AppEnv) *slog.Logger {
	logger := new(slog.Logger)

	switch env {
	case models.AppEnvLocal:
		opts := handlers.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		logger = slog.New(opts.NewPrettyHandler(os.Stdout))
	case models.AppEnvDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case models.AppEnvProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
