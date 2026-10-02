package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rqcoon/shimmer/api"
	"github.com/rqcoon/shimmer/cache"
	"github.com/rqcoon/shimmer/renderer"
	"github.com/rqcoon/shimmer/renderer/injector"
)

func main() {
	sourcedir := flag.String("dir", "./page", "Source directory for .md files")
	port := flag.String("port", "8080", "HTTP server port")
	flag.Parse()

	c := cache.New()
	r := renderer.New(*sourcedir)

	injector, err := injector.New()
	if err != nil {
		log.Fatalf("failed to create injector: %v", err)
	}

	server := api.NewServer(c, r, injector)

	if err := server.Start(); err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	addr := ":" + *port
	httpServer := &http.Server{
		Addr:    addr,
		Handler: server,
	}

	go func() {
		log.Printf("shimmer listening on %s\n", addr)

		if err := httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	<-signals

	log.Printf("caught signal interrupt\n")
	server.Stop()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	log.Println("shimmer stopped")
}
