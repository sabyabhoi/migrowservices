package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sabyabhoi/microservices/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Println("Server started on port 9090")
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			l.Fatalf("HTTP Server Error: %v", err)
		}
		l.Println("Stopped listening for new connections")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutDownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10 * time.Second)
	defer shutdownRelease()

	if err := s.Shutdown(shutDownCtx); err != nil {
		l.Fatalf("HTTP shutdown error: %v", err)
	}
	l.Println("Graceful shutdown complete")
}
