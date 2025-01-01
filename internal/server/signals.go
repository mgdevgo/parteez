package server

import (
	"os"
	"os/signal"
	"syscall"

	"log"
)

func (s *Server) interceptSignals() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		for {
			select {
			case sig := <-c:
				log.Printf("Intercepted %q signal", sig)
				switch sig {
				case syscall.SIGINT:
					s.Shutdown()
					s.WaitForShutdown()
					os.Exit(0)
					// case syscall.SIGTERM:
					// 	s.Shutdown()
					// 	s.WaitForShutdown()
					// 	os.Exit(1)
				}
			case <-s.quit:
				return
			}
		}
	}()
}
