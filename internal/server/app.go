package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	"parteez/internal/domain/shared"
	"parteez/internal/telegram"
)

type Server struct {
	http             *fiber.App
	storage          *pgxpool.Pool
	backgroundJobs   []string
	bot              *telegram.Bot
	logger           *slog.Logger
	commands         []*cobra.Command
	handlers         *Handlers
	quit             chan struct{}
	shutdownComplete chan struct{}
}

func New() (*Server, error) {
	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("parteez v%s", VERSION),
		// DisableStartupMessage: true,
		// ReadTimeout:  config.HTTPServer.Timeout * time.Second,
		// WriteTimeout: config.HTTPServer.Timeout * time.Second,
		// IdleTimeout:  config.HTTPServer.IdleTimeout * time.Second,
		ErrorHandler: shared.NewErrorHandler(),
	})

	app.Use(logger.New())

	limiterConfig := limiter.Config{
		LimitReached: func(ctx *fiber.Ctx) error {
			return &shared.Error{
				Status: fiber.StatusTooManyRequests,
				Code:   string(shared.ErrorCodeRateLimit),
			}
		},
	}
	app.Use(limiter.New(limiterConfig))

	app.Static("/telegram-mini-app", "./web/mini-app/build")

	server := &Server{
		http:             app,
		quit:             make(chan struct{}),
		shutdownComplete: make(chan struct{}),
		commands: []*cobra.Command{
			parteezCommand, // root command
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

var serverContextKey contextKey = "parteez-server"

func (server *Server) Start(args []string) error {
	parteez := server.commands[0]

	for _, cmd := range server.commands[1:] {
		parteez.AddCommand(cmd)
	}

	parteez.SetArgs(args)

	ctx := context.WithValue(context.Background(), serverContextKey, server)

	return parteez.ExecuteContext(ctx)
}

func (server *Server) Shutdown() {

	// Sguttdown http server
	server.http.Shutdown()

	// Close connections to storage
	server.storage.Close()

	// Release go routines that wait on that channel
	close(server.quit)

	// Do stuff

	// Notify that the shutdown is complete
	close(server.shutdownComplete)
}

// WaitForShutdown will block until the server has been fully shutdown.
func (server *Server) WaitForShutdown() {
	<-server.shutdownComplete
}
