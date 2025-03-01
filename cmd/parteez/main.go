package main

import (
	"fmt"
	"parteez/internal/server"
	"os"
)

func main() {
	exe := "iditusi-server"

	// Create the server with appropriate options.
	app, err := server.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s", exe, err)
		os.Exit(1)
	}

	// Start things up. Block here until done.
	if err := app.Start(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s", exe, err)
		os.Exit(1)
	}

	app.WaitForShutdown()
}
