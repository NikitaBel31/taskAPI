package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taskapi/internal/app"
	"time"
)

func main() {
	c := app.NewContainer()
	defer c.Logger.Stop()

	srv := &http.Server{
		Addr:              c.Config.HTTPPort,
		Handler:           c.Router.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("HTTP server listening on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Config.ShutdownTime)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("stopped")
}
