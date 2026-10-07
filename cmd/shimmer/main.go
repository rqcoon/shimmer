package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rqcoon/shimmer"
)

func main() {
	contentDir := flag.String("dir", "./page", "Source directory for Markdown files")
	templateDir := flag.String("templates", "./templates", "Template directory")
	port := flag.String("port", "8080", "HTTP server port")

	flag.Parse()

	s, err := shimmer.New(shimmer.Config{
		ContentDir:  *contentDir,
		TemplateDir: *templateDir,
		WatchRate:   500 * time.Millisecond,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: s.Handler(),
	}

	go func() {
		log.Printf("shimmer listening on %s", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	log.Println("shutting down shimmer")

	s.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}
}
