package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

func main() {
	iditusiCommand := &cobra.Command{
		Use: "iditusi",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	iditusiCommand.AddCommand(
		NewServeCommand(),
		NewRoutesCommand(),
		NewMigrateCommand(),
	)

	if err := iditusiCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
