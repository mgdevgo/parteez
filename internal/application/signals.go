package application

import (
	"os"
	"os/signal"
	"syscall"

	"log"
)

func (app *Application) interceptSignals() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		for {
			select {
			case sig := <-c:
				log.Printf("Intercepted %q signal", sig)
				switch sig {
				case syscall.SIGINT:
					app.Shutdown()
					app.WaitForShutdown()
					os.Exit(0)
					// case syscall.SIGTERM:
					// 	s.Shutdown()
					// 	s.WaitForShutdown()
					// 	os.Exit(1)
				}
			case <-app.quit:
				return
			}
		}
	}()
}
