package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"parteez/internal/app"
	"parteez/internal/config"
	"parteez/pkg/logger"
)

func main() {
	root, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	//done := make(chan os.Signal, 1)
	//signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	//
	//<-done
	conf := config.New()
	log := logger.New(conf.AppEnv)

	log.Info("Starting application",
		slog.String("Name", "Parteez"),
		slog.String("Version", "0.1.0"),
		slog.String("ENV", conf.AppEnv),
	)

	db, err := pgxpool.New(context.TODO(), conf.DatabaseURL)
	if err != nil {
		log.Error("Can not initialize database connection pool",
			slog.String("URL", conf.DatabaseURL),
			slog.String("Error", err.Error()),
		)
		os.Exit(1)
	}
	if err := db.Ping(context.TODO()); err != nil {
		log.Error("Failed to ping database",
			slog.String("URL", conf.DatabaseURL),
			slog.String("Error", err.Error()),
		)
		os.Exit(1)
	}

	log.Info("Database connection OK")

	application := app.New(conf, db, log)

	g, run := errgroup.WithContext(root)
	g.Go(func() error {
		log.Info("HTTPServer started accepting connections", slog.String("Address", conf.HTTPServer.Address))

		if err := application.HTTPServer.Listen(conf.HTTPServer.Address); err != nil {
			log.Error("HTTPServer failed to listen", slog.String("Error", err.Error()))
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-run.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := application.HTTPServer.ShutdownWithContext(ctx); err != nil {
			log.Error("HTTPServer stopped with error", slog.String("Error", err.Error()))
			return err
		}
		log.Info("HTTPServer stopped")
		return nil
	})

	g.Go(func() error {
		<-run.Done()
		application.Storage.Close()
		log.Info("Database connections closed")
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Error("Application stopped with error", slog.String("Error", err.Error()))
	} else {
		log.Info("Application stopped")
	}
}
