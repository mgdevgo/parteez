package app

import (
	"github.com/gofiber/fiber/v2"
	console "github.com/urfave/cli/v2"
)

type Application struct {
	app         *fiber.App
	middlewares []fiber.Handler
	commands    []*console.Command
}

func (a *Application) Start(args []string) error {
	command := &console.App{
		Name:     "iditusi",
		Commands: a.commands,
	}

	return command.Run(args)
}

func (a *Application) Stop() error
func (a *Application) Routes() string
