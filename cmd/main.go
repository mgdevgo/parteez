package main

import (
	"context"
	"fmt"
	"os"

	"iditusi/internal/application"
)

const version = "0.1.0"

func main() {
	ctx := context.Background()

	if err := application.Run(ctx, os.Args); err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}
}
