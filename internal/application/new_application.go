package application

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

type Middleware interface {
}

type AnyCommand interface {
	// Text that will be displayed when --help is passed.
	Help() string
	// Runs the command
	Run() error
}

type App struct {
	middlewares []fiber.Handler
	commands    []*cobra.Command
}
