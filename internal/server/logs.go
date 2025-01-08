package server

import (
	"iditusi/internal/log/handlers"
	"log/slog"
	"os"
)

// Logger - returns current logger
func (s *Server) Logger() *slog.Logger {
	// s.logging.Lock()
	// defer s.logging.Unlock()
	// return s.logging.logger
	return s.logger
}

func (s *Server) WithDebug(logger *slog.Logger) *Server {
	if logger != nil {
		s.logger = logger
		return s
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

	return s
}
