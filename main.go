package main

import (
	"fmt"
	"iditusi/internal/server"
	"os"
)

func main() {
	exe := "iditusi-server"

	// Create the server with appropriate options.
	s, err := server.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: %s", exe, err))
		os.Exit(1)
	}

	// Start things up. Block here until done.
	if err := s.Start(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: %s", exe, err))
		os.Exit(1)
	}

	s.WaitForShutdown()
}
