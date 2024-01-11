package main

import (
	"context"
	"fmt"
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

	conf := config.New()
	log := logger.New(conf.AppEnv)

	pool, err := pgxpool.New(root, conf.DatabaseURL)
	if err != nil {
		panic(fmt.Errorf("pgxpool.New: %w", err))
	}

	ctx, cancel := context.WithTimeout(root, time.Second*2)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Errorf("ping: %w", err))
	}

	log.Info("Connected to database", slog.String("connection_url", conf.DatabaseURL))

	application := app.New(conf, pool, log)

	g, run := errgroup.WithContext(root)
	g.Go(func() error {
		if err := application.HTTPServer.Listen(conf.HTTPServer.Address); err != nil {
			log.Error("HTTP server listen error", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-run.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := application.HTTPServer.ShutdownWithContext(ctx); err != nil {
			log.Error("Failed to shutdown HTTP server", slog.String("error", err.Error()))
			return err
		}
		log.Info("HTTP server stopped")
		return nil
	})

	g.Go(func() error {
		<-run.Done()
		application.Storage.Close()
		log.Info("Database connections closed")
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Error("Application stopped with error", slog.String("error", err.Error()))
	}

	log.Info("Application stopped")
}
