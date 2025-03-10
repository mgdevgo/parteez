package application

import (
	"log/slog"
	"os"

	"parteez/pkg/log/handlers"
)

// Logger - returns current logger
func (app *Application) Logger() *slog.Logger {
	// s.logging.Lock()
	// defer s.logging.Unlock()
	// return s.logging.logger
	return app.logger
}

func (app *Application) WithDebug(logger *slog.Logger) *Application {
	if logger != nil {
		app.logger = logger
		return app
	}

	opts := handlers.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	logger = slog.New(opts.NewPrettyHandler(os.Stdout))

	// logger = slog.New(
	// 	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	// )

	// logger = slog.New(
	// 	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	// )

	return app
}
