package server

import "log/slog"

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

	s.logger = slog.Default()

	return s
}
