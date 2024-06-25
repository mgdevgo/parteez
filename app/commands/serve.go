package commands

import (
	"fmt"

	console "github.com/urfave/cli/v2"
)

var ServeCommand = &console.Command{
	Name:  "serve",
	Usage: "Begins serving the app over HTTP.",
	Flags: []console.Flag{
		&console.StringFlag{
			Name:    "hostname",
			Usage:   "Set the hostname the server will run on.",
			Aliases: []string{"H"},
		},
		&console.StringFlag{
			Name:    "port",
			Usage:   "Set the port the server will run on.",
			Aliases: []string{"p"},
		},
		&console.StringFlag{
			Name:    "bind",
			Usage:   "Convenience for setting hostname and port together.",
			Aliases: []string{"b"},
		},
	},
	Action: func(context *console.Context) error {
		fmt.Println(context.String("hostname"))
		// app := context.Context.Value("iditusi.application")
		return nil
	},
}
