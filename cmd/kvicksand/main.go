package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hailkomputer/kvicksand/internal/api"
	"github.com/hailkomputer/kvicksand/internal/cache"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      api.NewApiHandler(cache.NewCache()).Router,
	}

	go func() {
		log.Println("starting API server")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Println(err)
		}
		log.Println("gracefully stopping API Server")
	}()

	c := make(chan os.Signal, 1)

	// Shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	return srv.Shutdown(ctx)
}
