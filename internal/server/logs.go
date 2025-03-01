package server

import (
	"parteez/pkg/log/handlers"
	"log/slog"
	"os"
)

// Logger - returns current logger
func (app *Server) Logger() *slog.Logger {
	// s.logging.Lock()
	// defer s.logging.Unlock()
	// return s.logging.logger
	return app.logger
}

func (app *Server) WithDebug(logger *slog.Logger) *Server {
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
