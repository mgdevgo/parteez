package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"parteez/internal/app"
	"parteez/internal/config"
)

func main() {
	conf := config.MustLoad()

	main, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	//done := make(chan os.Signal, 1)
	//signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	//
	//<-done

	//logger := logger.New()

	application, _ := app.New(main, conf)

	g, run := errgroup.WithContext(main)
	g.Go(func() error {
		const op = "HTTPServer.Listen"
		log.Printf("API server listening at: %s \n", conf.HTTPServer.Address)

		if err := application.HTTPServer.Listen("123"); err != nil {
			log.Printf("%s: %s", op, err.Error())
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	g.Go(func() error {
		<-run.Done()

		// TODO: move timeout to config
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := application.HTTPServer.ShutdownWithContext(ctx); err != nil {
			const op = "HTTPServer.ShutdownWithContext"
			log.Printf("%s: %s", op, err.Error())

			return fmt.Errorf("%s: %w", op, err)
		}

		log.Println("HTTPServer stopped")
		return nil
	})
	g.Go(func() error {
		<-run.Done()
		application.Storage.Close()
		log.Println("Storage closed")

		return nil
	})

	if err := g.Wait(); err != nil {
		log.Printf("Application stopped with err: %v\n", err.Error())
	} else {
		log.Println("Application stopped")
	}

}
