package application

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	"parteez/internal/domain/shared/handler"
	"parteez/internal/telegram"
	"parteez/pkg/environment"
)

type Application struct {
	env              environment.Environment
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

func New(env ...environment.Environment) (*Application, error) {
	appenv := environment.Development
	if len(env) > 0 {
		appenv = env[0]
	}

	errorHandler := handler.NewErrorHandler(nil)

	fiberConfig := fiber.Config{
		AppName: fmt.Sprintf("parteez v%s", VERSION),
		// DisableStartupMessage: true,
		// ReadTimeout:  config.HTTPServer.Timeout * time.Second,
		// WriteTimeout: config.HTTPServer.Timeout * time.Second,
		// IdleTimeout:  config.HTTPServer.IdleTimeout * time.Second,
		ErrorHandler: errorHandler.Handle,
	}
	fiberApp := fiber.New(fiberConfig)

	rateLimitConfig := limiter.Config{
		LimitReached: func(ctx *fiber.Ctx) error {
			return &handler.Error{
				Status: fiber.StatusTooManyRequests,
				Code:   string(handler.ErrorCodeRateLimit),
			}
		},
	}

	fiberApp.Use(
		limiter.New(rateLimitConfig),
		logger.New(),
	)

	fiberApp.Static("/telegram-mini-app", "./web/mini-app/build")

	application := &Application{
		env:              appenv,
		http:             fiberApp,
		quit:             make(chan struct{}),
		shutdownComplete: make(chan struct{}),
		backgroundJobs:   make([]string, 0),
		commands:         make([]*cobra.Command, 0),
	}

	application.configureRoutes()
	application.configureCommands()

	return application, nil
}

func (app *Application) Start(args []string) error {
	combinedCommands := &cobra.Command{
		Use: "parteez",
	}

	for _, command := range app.commands {
		combinedCommands.AddCommand(command)
	}

	combinedCommands.SetArgs(args)

	ctx := contextWithApplication(combinedCommands.Context(), app)

	return combinedCommands.ExecuteContext(ctx)
}

func (app *Application) Execute() error {
	go func() {
		if err := app.http.Listen(":8080"); err != nil {
			app.logger.Error("Failed to start http server", "error", err)
			close(app.quit)
		}
	}()

	app.interceptSignals()

	app.WaitForShutdown()

	return nil
}

func (server *Application) Shutdown() {

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
func (server *Application) WaitForShutdown() {
	<-server.shutdownComplete
}
