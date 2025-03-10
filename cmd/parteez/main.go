package main

import (
	"fmt"
	"os"

	"parteez/internal/application"
)

func main() {
	// Create the server with appropriate options.
	app, err := application.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	// Start things up. Block here until done.
	if err := app.Start(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}
