package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewMigrateCommand() *cobra.Command {
	migrateCommand := &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("migrate command not implemented")
		},
	}
	return migrateCommand
}
