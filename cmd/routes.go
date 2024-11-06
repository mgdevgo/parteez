package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRoutesCommand() *cobra.Command {
	routesCommand := &cobra.Command{
		Use: "routes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("routes command not implemented")
		},
	}
	return routesCommand
}
