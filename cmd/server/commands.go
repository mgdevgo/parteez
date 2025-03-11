package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// const applicationContextKey = "application"

// func contextWithApplication(ctx context.Context, application *Application) context.Context {
// 	return context.WithValue(ctx, applicationContextKey, application)
// }

// func applicationFromContext(ctx context.Context) (*Application, error) {
// 	application, ok := ctx.Value(applicationContextKey).(*Application)
// 	if !ok {
// 		return nil, errors.New("application not set on context")
// 	}
// 	return application, nil
// }

func parteezCommand() *cobra.Command {
	return &cobra.Command{
		Use: "parteez",
		RunE: func(cmd *cobra.Command, args []string) error {
			// ctx := cmd.Context()
			// app, err := applicationFromContext(ctx)
			// if err != nil {
			// 	return err
			// }
			// return app.Start(args)
			return nil
		},
	}
}

func migrateCommand() *cobra.Command {
	return &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("migrate command not implemented")
		},
	}
}

func routesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "routes",
		Short: "Displays all registered routes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// ctx := cmd.Context()
			// app, err := applicationFromContext(ctx)
			// if err != nil {
			// 	return err
			// }
			// routes := app.http.GetRoutes()
			// for _, route := range routes {
			// 	fmt.Println(route.Path, route.Method)
			// }
			return nil
		},
	}
}

func serveCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Begins serving the app over HTTP.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// ctx := cmd.Context()
			// app, err := applicationFromContext(ctx)
			// if err != nil {
			// 	return err
			// }
			// return app.Start(args)
			return nil
		},
	}

	command.Flags().String("bind", "0.0.0.0:8080", "Convenience for setting hostname and port together.")
	command.Flags().StringP("hostname", "h", "0.0.0.0", "Set the hostname the server will run on.")
	command.Flags().StringP("port", "p", "8080", "Set the port the server will run on.")

	return command
}
