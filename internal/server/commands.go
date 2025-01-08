package server

import (
	"errors"
	"fmt"
	"iditusi/internal/storage/postgres"

	"github.com/spf13/cobra"
)

var iditusiCommand = &cobra.Command{
	Use: "iditusi",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		server, ok := ctx.Value(serverContextKey).(*Server)
		if !ok {
			return errors.New("can't get server reference from context")
		}
		// Do Stuff Here
		server.WaitForShutdown()
		return nil
	},
}

var migrateCommand = &cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrate command not implemented")
	},
}

var routesCommand = &cobra.Command{
	Use: "routes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("routes command not implemented")
	},
}

var serveCommand = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := postgres.New("")
		cobra.CheckErr(err)
		_ = db
		return nil
	},
}

func init() {
	serveCommand.Flags().String("bind", "0.0.0.0:8080", "Convenience for setting hostname and port together.")
	serveCommand.Flags().StringP("hostname", "h", "0.0.0.0", "Set the hostname the server will run on.")
	serveCommand.Flags().StringP("port", "p", "8080", "Set the port the server will run on.")
}
