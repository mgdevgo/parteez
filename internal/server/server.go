package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"

	"iditusi/internal/storage/postgres"
	"iditusi/internal/telegram"
)

type Server struct {
	httpServer     *fiber.App
	storage        *postgres.Storage
	backgroundJobs []string
	bot            *telegram.Bot
	logger         *slog.Logger
	commands       []*cobra.Command

	quit             chan struct{}
	shutdownComplete chan struct{}
}

func New() (*Server, error) {
	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("iditusi v%s", VERSION),
		// DisableStartupMessage: true,
		// ReadTimeout:  config.HTTPServer.Timeout * time.Second,
		// WriteTimeout: config.HTTPServer.Timeout * time.Second,
		// IdleTimeout:  config.HTTPServer.IdleTimeout * time.Second,
		// ErrorHandler: api.ErrorHandler(),
	})

	app.Static("/telegram-mini-app", "./web/mini-app/build")

	server := &Server{
		httpServer:       app,
		quit:             make(chan struct{}),
		shutdownComplete: make(chan struct{}),
		commands: []*cobra.Command{
			iditusiCommand, // root command
			serveCommand,
			routesCommand,
			migrateCommand,
		},
	}

	// Start signal intercepter
	server.interceptSignals()

	return server, nil
}

type contextKey string

var serverContextKey contextKey = "iditusi-server"

func (s *Server) Start(args []string) error {
	iditusi := s.commands[0]

	for _, cmd := range s.commands[1:] {
		iditusi.AddCommand(cmd)
	}

	iditusi.SetArgs(args)

	ctx := context.WithValue(context.Background(), serverContextKey, s)

	return iditusi.ExecuteContext(ctx)
}

func (s *Server) Shutdown() {

	// Sguttdown http server
	s.httpServer.Shutdown()

	// Close connections to storage
	s.storage.Close()

	// Release go routines that wait on that channel
	close(s.quit)

	// Do stuff

	// Notify that the shutdown is complete
	close(s.shutdownComplete)
}

// WaitForShutdown will block until the server has been fully shutdown.
func (s *Server) WaitForShutdown() {
	<-s.shutdownComplete
}
