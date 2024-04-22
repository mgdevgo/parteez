package main

import (
	"context"
	"fmt"
	"os"

	"iditusi/pkg/app"
)

func main() {
	ctx := context.Background()

	if err := app.Run(ctx, os.Args); err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}
}
